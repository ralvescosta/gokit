package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
)

// New(...) create a new instance for Imessaging
//
// New(...) connect to the RabbitMQ broker and stablish a channel
func NewClient(cfg *env.RabbitMQConfigs, logger logging.Logger) Messaging {
	rb := &messagingImpl{
		logger: logger,
		config: cfg,
		tracer: otel.Tracer("rabbitmq"),
	}

	logger.Debug(Message("connecting to rabbitmq..."))
	conn, err := dial(cfg)
	if err != nil {
		logger.Error(Message("failure to connect to the broker"), zap.Error(err))
		rb.Err = ErrorAMQPConnection
		return rb
	}
	logger.Debug(Message("connected to rabbitmq"))

	rb.conn = conn

	logger.Debug(Message("creating amqp channel..."))
	ch, err := conn.Channel()
	if err != nil {
		logger.Error(Message("failure to establish the channel"), zap.Error(err))
		rb.Err = ErrorAMQPChannel
		return rb
	}
	logger.Debug(Message("created amqp channel"))

	rb.channel = ch

	return rb
}

var dial = func(cfg *env.RabbitMQConfigs) (AMQPConnection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.VHost, cfg.Port))
}

func (m *messagingImpl) InstallTopology(tplogy Topology) (Messaging, error) {
	tp := tplogy.(*topologyImpl)

	if m.Err != nil {
		return nil, m.Err
	}

	for _, opts := range tp.exchanges {
		m.logger.Debug(Message("declaring exchanges..."))
		if err := m.installExchange(opts); err != nil {
			m.logger.Error(Message("declare exchange err"), zap.Error(err))
			return nil, err
		}
		m.logger.Debug(Message("exchanges declared"))
	}

	for _, opts := range tp.queues {
		m.logger.Debug(Message("declaring queues..."))
		if err := m.installQueues(opts); err != nil {
			m.logger.Error(Message("declare queue err"), zap.Error(err))
			return nil, err
		}
		m.logger.Debug(Message("queues declared"))
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

func (m *messagingImpl) Channel() AMQPChannel {
	return m.channel
}

func (m *messagingImpl) installExchange(opt *ExchangeOpts) error {
	err := m.channel.ExchangeDeclare(opt.name, string(opt.kind), true, false, false, false, nil)

	if err != nil {
		return err
	}

	return nil
}

func (m *messagingImpl) installQueues(opts *QueueOpts) error {
	var amqpDlqDeclarationOpts amqp.Table

	if opts.retry != nil {
		m.logger.Debug(Message("declaring retry queue..."))

		_, err := m.channel.QueueDeclare(opts.RetryName(), true, false, false, false, amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": opts.name,
			"x-message-ttl":             opts.retry.DelayBetween.Milliseconds(),
		})

		if err != nil {
			return err
		}

		amqpDlqDeclarationOpts = amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": opts.RetryName(),
		}
		m.logger.Debug(Message("retry queue declared"))
	}

	if amqpDlqDeclarationOpts == nil && opts.withDeadLatter {
		amqpDlqDeclarationOpts = amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": opts.DqlName(),
		}
	}

	if opts.withDeadLatter {
		m.logger.Debug(Message("declaring dlq queue..."))
		_, err := m.channel.QueueDeclare(opts.DqlName(), true, false, false, false, amqpDlqDeclarationOpts)

		if err != nil {
			return err
		}
		m.logger.Debug(Message("dlq queue declared"))
	}

	_, err := m.channel.QueueDeclare(opts.name, true, false, false, false, amqpDlqDeclarationOpts)

	if err != nil {
		return err
	}

	for _, biding := range opts.bindings {
		m.logger.Debug(Message("binding queue..."))
		err := m.channel.QueueBind(opts.name, biding.routingKey, biding.exchange, false, nil)

		if err != nil {
			return err
		}
		m.logger.Debug(Message("queue bonded"))
	}

	return nil
}
