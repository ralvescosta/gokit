package rabbitmq

import (
	"github.com/ralvescosta/gokit/logging"
	"github.com/streadway/amqp"
)

type (
	Topology struct {
		logger           logging.Logger
		channel          AMQPChannel
		queues           map[string]*QueueDefinition
		exchanges        []*ExchangeDefinition
		exchangeBindings []*ExchangeBindingDefinition
		queueBindings    map[string]*QueueBindingDefinition
	}
)

func NewTopology(l logging.Logger) *Topology {
	return &Topology{logger: l}
}

func (t *Topology) Channel(c AMQPChannel) *Topology {
	t.channel = c
	return t
}

func (t *Topology) Queue(q *QueueDefinition) *Topology {
	t.queues[q.name] = q
	return t
}

func (t *Topology) Queues(queues []*QueueDefinition) *Topology {
	for _, q := range queues {
		t.queues[q.name] = q
	}

	return t
}

func (t *Topology) Exchange(e *ExchangeDefinition) *Topology {
	t.exchanges = append(t.exchanges, e)
	return t
}

func (t *Topology) Exchanges(e []*ExchangeDefinition) *Topology {
	t.exchanges = append(t.exchanges, e...)
	return t
}

func (t *Topology) ExchangeBinding(b *ExchangeBindingDefinition) *Topology {
	t.exchangeBindings = append(t.exchangeBindings, b)
	return t
}

func (t *Topology) QueueBinding(b *QueueBindingDefinition) *Topology {
	t.queueBindings[b.queue] = b
	return t
}

func (t *Topology) Apply() error {
	if t.channel == nil {
		return NullableChannel
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

func (t *Topology) declareExchanges() error {
	t.logger.Debug(LogMessage("declaring exchanges..."))

	for _, exch := range t.exchanges {
		if err := t.channel.ExchangeDeclare(exch.name, exch.kind.String(), exch.durable, exch.delete, false, false, exch.params); err != nil {
			return err
		}
	}

	t.logger.Debug(LogMessage("exchanges declared"))

	return nil
}

func (t *Topology) declareQueues() error {
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

func (t *Topology) bindQueues() error {
	t.logger.Debug(LogMessage("binding queues..."))

	for _, bind := range t.queueBindings {
		if err := t.channel.QueueBind(bind.queue, bind.routingKey, bind.exchange, false, bind.args); err != nil {
			return err
		}
	}

	t.logger.Debug(LogMessage("queue bonded"))

	return nil
}

func (t *Topology) bindExchanges() error {
	t.logger.Debug(LogMessage("binding exchanges..."))

	for _, bind := range t.exchangeBindings {
		if err := t.channel.ExchangeBind(bind.destination, bind.routingKey, bind.source, false, bind.args); err != nil {
			return err
		}
	}

	t.logger.Debug(LogMessage("exchanges bonded"))

	return nil
}
