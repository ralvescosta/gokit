package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/errors"
	"github.com/ralvescosta/gokit/logging"
)

// New(...) create a new instance for Imessaging
//
// New(...) connect to the RabbitMQ broker and stablish a channel
func New(cfg *env.Configs, logger logging.ILogger) Messaging {
	rb := &messaging{
		logger: logger,
		config: cfg,
		tracer: otel.Tracer("rabbitmq"),
	}

	logger.Debug(LogMessage("connecting to rabbitmq..."))
	conn, err := dial(cfg)
	if err != nil {
		logger.Error(LogMessage("failure to connect to the broker"), logging.ErrorField(err))
		rb.Err = errors.ErrorAMQPConnection
		return rb
	}
	logger.Debug(LogMessage("connected to rabbitmq"))

	rb.conn = conn

	logger.Debug(LogMessage("creating amqp channel..."))
	ch, err := conn.Channel()
	if err != nil {
		logger.Error(LogMessage("failure to establish the channel"), logging.ErrorField(err))
		rb.Err = errors.ErrorAMQPChannel
		return rb
	}
	logger.Debug(LogMessage("created amqp channel"))

	rb.channel = ch

	return rb
}

var dial = func(cfg *env.Configs) (AMQPConnection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.RABBIT_USER, cfg.RABBIT_PASSWORD, cfg.RABBIT_VHOST, cfg.RABBIT_PORT))
}

func (m *messaging) InstallTopology(topology *topology) (Messaging, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	for _, opts := range topology.exchanges {
		m.logger.Debug(LogMessage("declaring exchanges..."))
		if err := m.installExchange(opts); err != nil {
			m.logger.Error(LogMessage("declare exchange err"), logging.ErrorField(err))
			return nil, err
		}
		m.logger.Debug(LogMessage("exchanges declared"))
	}

	for _, opts := range topology.queues {
		m.logger.Debug(LogMessage("declaring queues..."))
		if err := m.installQueues(opts); err != nil {
			m.logger.Error(LogMessage("declare queue err"), logging.ErrorField(err))
			return nil, err
		}
		m.logger.Debug(LogMessage("queues declared"))
	}

	return m, m.Err
}

// func (m *RabbitMQMessaging) Publisher(exchange, routingKey string, msg any, opts *PublishOpts) error {
// 	byt, err := json.Marshal(msg)
// 	if err != nil {
// 		m.logger.Error(LogMessage("publisher marshal"), logging.ErrorField(err))
// 		return err
// 	}

// 	if opts == nil {
// 		opts = m.newPubOpts(fmt.Sprintf("%T", msg))
// 	}

// 	return m.channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
// 		Headers: amqp.Table{
// 			AMQPHeaderNumberOfRetry: opts.Count,
// 			AMQPHeaderTraceparent:   opts.Traceparent,
// 			AMQPHeaderDelay:         opts.Delay.Milliseconds(),
// 		},
// 		Type:        opts.Type,
// 		ContentType: JsonContentType,
// 		MessageId:   opts.MessageId,
// 		UserId:      m.config.RABBIT_USER,
// 		AppId:       m.config.APP_NAME,
// 		Body:        byt,
// 	})
// }

func (m *messaging) Channel() AMQPChannel {
	return m.channel
}

func (m *messaging) installExchange(opt *ExchangeOpts) error {
	err := m.channel.ExchangeDeclare(opt.name, string(opt.kind), true, false, false, false, nil)

	if err != nil {
		return err
	}

	return nil
}

func (m *messaging) installQueues(opts *QueueOpts) error {
	var amqpDlqDeclarationOpts amqp.Table

	if opts.retry != nil {
		m.logger.Debug(LogMessage("declaring retry queue..."))
		retryQueueName := fmt.Sprintf("%s-retry", opts.name)

		_, err := m.channel.QueueDeclare(retryQueueName, true, false, false, false, amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": opts.name,
			"x-message-ttl":             opts.retry.DelayBetween.Milliseconds(),
		})

		if err != nil {
			return err
		}

		amqpDlqDeclarationOpts = amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": retryQueueName,
		}
		m.logger.Debug(LogMessage("retry queue declared"))
	}

	dlqQueueName := fmt.Sprintf("%s-dlq", opts.name)
	if amqpDlqDeclarationOpts == nil && opts.withDeadLatter {
		amqpDlqDeclarationOpts = amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": dlqQueueName,
		}
	}

	if opts.withDeadLatter {
		m.logger.Debug(LogMessage("declaring dlq queue..."))
		_, err := m.channel.QueueDeclare(dlqQueueName, true, false, false, false, amqpDlqDeclarationOpts)

		if err != nil {
			return err
		}
		m.logger.Debug(LogMessage("dlq queue declared"))
	}

	_, err := m.channel.QueueDeclare(opts.name, true, false, false, false, amqpDlqDeclarationOpts)

	if err != nil {
		return err
	}

	for _, biding := range opts.bindings {
		m.logger.Debug(LogMessage("binding queue..."))
		err := m.channel.QueueBind(opts.name, biding.routingKey, biding.exchange, false, nil)

		if err != nil {
			return err
		}
		m.logger.Debug(LogMessage("queue bonded"))
	}

	return nil
}
