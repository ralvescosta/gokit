package rabbitmq

import (
	"github.com/ralvescosta/gokit/logging"
	"github.com/streadway/amqp"
)

type (
	Topology interface {
		Channel(c AMQPChannel) Topology
		Queue(q *QueueDefinition) Topology
		Queues(queues []*QueueDefinition) Topology
		Exchange(e *ExchangeDefinition) Topology
		Exchanges(e []*ExchangeDefinition) Topology
		ExchangeBinding(b *ExchangeBindingDefinition) Topology
		QueueBinding(b *QueueBindingDefinition) Topology
		GetQueuesDefinition() map[string]*QueueDefinition
		GetQueueDefinition(queueName string) (*QueueDefinition, error)
		Apply() error
	}

	topology struct {
		logger           logging.Logger
		channel          AMQPChannel
		queues           map[string]*QueueDefinition
		queuesBinding    map[string]*QueueBindingDefinition
		exchanges        []*ExchangeDefinition
		exchangesBinding []*ExchangeBindingDefinition
	}
)

func NewTopology(l logging.Logger) *topology {
	return &topology{logger: l}
}

func (t *topology) Channel(c AMQPChannel) *topology {
	t.channel = c
	return t
}

func (t *topology) Queue(q *QueueDefinition) *topology {
	t.queues[q.name] = q
	return t
}

func (t *topology) Queues(queues []*QueueDefinition) *topology {
	for _, q := range queues {
		t.queues[q.name] = q
	}

	return t
}

func (t *topology) GetQueuesDefinition(queueName string) map[string]*QueueDefinition {
	return t.queues
}

func (t *topology) GetQueueDefinition(queueName string) (*QueueDefinition, error) {
	if d, ok := t.queues[queueName]; ok {
		return d, nil
	}

	return nil, NotFoundQueueDefinitionError
}

func (t *topology) Exchange(e *ExchangeDefinition) *topology {
	t.exchanges = append(t.exchanges, e)
	return t
}

func (t *topology) Exchanges(e []*ExchangeDefinition) *topology {
	t.exchanges = append(t.exchanges, e...)
	return t
}

func (t *topology) ExchangeBinding(b *ExchangeBindingDefinition) *topology {
	t.exchangesBinding = append(t.exchangesBinding, b)
	return t
}

func (t *topology) QueueBinding(b *QueueBindingDefinition) *topology {
	t.queuesBinding[b.queue] = b
	return t
}

func (t *topology) Apply() error {
	if t.channel == nil {
		return NullableChannelError
	}

	if err := t.declareExchanges(); err != nil {
		return err
	}

	if err := t.declareQueues(); err != nil {
		return err
	}

	if err := t.bindQueues(); err != nil {
		return err
	}

	return t.bindExchanges()
}

func (t *topology) declareExchanges() error {
	t.logger.Debug(LogMessage("declaring exchanges..."))

	for _, exch := range t.exchanges {
		if err := t.channel.ExchangeDeclare(exch.name, exch.kind.String(), exch.durable, exch.delete, false, false, exch.params); err != nil {
			return err
		}
	}

	t.logger.Debug(LogMessage("exchanges declared"))

	return nil
}

func (t *topology) declareQueues() error {
	for _, queue := range t.queues {
		if queue.withRetry {
			t.logger.Debug(LogMessage("declaring retry queue..."))

			//queue.RetryName(), true, false, false, false, amqpDlqDeclarationOpts
			if _, err := t.channel.QueueDeclare(queue.RetryName(), queue.durable, queue.delete, queue.exclusive, false, amqp.Table{
				"x-dead-letter-exchange":    "",
				"x-dead-letter-routing-key": queue.name,
				"x-message-ttl":             queue.retryTTL.Milliseconds(),
			}); err != nil {
				return err
			}

			t.logger.Debug(LogMessage("retry queue declared"))
		}

		var amqpDlqDeclarationOpts amqp.Table
		if queue.withDLQ && queue.withRetry {
			amqpDlqDeclarationOpts = amqp.Table{
				"x-dead-letter-exchange":    "",
				"x-dead-letter-routing-key": queue.RetryName(),
			}
		}

		if queue.withDLQ && !queue.withRetry {
			amqpDlqDeclarationOpts = amqp.Table{
				"x-dead-letter-exchange":    "",
				"x-dead-letter-routing-key": queue.DLQName(),
			}
		}

		if queue.withDLQ {
			t.logger.Debug(LogMessage("declaring dlq queue..."))

			//queue.DLQName(), true, false, false, false, amqpDlqDeclarationOpts
			if _, err := t.channel.QueueDeclare(queue.DLQName(), queue.durable, queue.delete, queue.exclusive, false, amqpDlqDeclarationOpts); err != nil {
				return err
			}

			t.logger.Debug(LogMessage("dlq queue declared"))

			if _, err := t.channel.QueueDeclare(queue.name, queue.durable, queue.delete, queue.exclusive, false, amqpDlqDeclarationOpts); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *topology) bindQueues() error {
	t.logger.Debug(LogMessage("binding queues..."))

	for _, bind := range t.queuesBinding {
		if err := t.channel.QueueBind(bind.queue, bind.routingKey, bind.exchange, false, bind.args); err != nil {
			return err
		}
	}

	t.logger.Debug(LogMessage("queue bonded"))

	return nil
}

func (t *topology) bindExchanges() error {
	t.logger.Debug(LogMessage("binding exchanges..."))

	for _, bind := range t.exchangesBinding {
		if err := t.channel.ExchangeBind(bind.destination, bind.routingKey, bind.source, false, bind.args); err != nil {
			return err
		}
	}

	t.logger.Debug(LogMessage("exchanges bonded"))

	return nil
}
