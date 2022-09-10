package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/errors"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/otel/trace"
)

// New(...) create a new instance for IRabbitMQMessaging
//
// New(...) connect to the RabbitMQ broker and stablish a channel
func New(cfg *env.Configs, logger logging.ILogger) IRabbitMQMessaging {
	rb := &RabbitMQMessaging{
		logger:      logger,
		config:      cfg,
		dispatchers: []*Dispatcher{},
		tracer:      otel.Tracer("rabbitmq"),
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

func (m *RabbitMQMessaging) InstallTopology(topology *Topology) (IRabbitMQMessaging, error) {
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

func (m *RabbitMQMessaging) Publisher(exchange, routingKey string, msg any, opts *PublishOpts) error {
	byt, err := json.Marshal(msg)
	if err != nil {
		m.logger.Error(LogMessage("publisher marshal"), logging.ErrorField(err))
		return err
	}

	if opts == nil {
		opts = m.newPubOpts(fmt.Sprintf("%T", msg))
	}

	return m.channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		Headers: amqp.Table{
			AMQPHeaderNumberOfRetry: opts.Count,
			AMQPHeaderTraceparent:   opts.Traceparent,
			AMQPHeaderDelay:         opts.Delay.Milliseconds(),
		},
		Type:        opts.Type,
		ContentType: JsonContentType,
		MessageId:   opts.MessageId,
		UserId:      m.config.RABBIT_USER,
		AppId:       m.config.APP_NAME,
		Body:        byt,
	})
}

func (m *RabbitMQMessaging) RegisterDispatcher(queue string, handler ConsumerHandler, t any) error {
	if t == nil || queue == "" {
		return errors.ErrorAMQPRegisterDispatcher
	}

	var conf *Topology

	for _, v := range m.topologies {
		if v.Queue.Name == queue {
			conf = v
			break
		}
	}

	dispatch := &Dispatcher{
		Queue:         queue,
		Topology:      conf,
		Handler:       handler,
		MsgType:       fmt.Sprintf("%T", t),
		ReflectedType: reflect.New(reflect.TypeOf(t).Elem()),
	}

	m.dispatchers = append(m.dispatchers, dispatch)

	return nil
}

func (m *RabbitMQMessaging) Consume() error {
	if m.Err != nil {
		return m.Err
	}

	m.shotdown = make(chan error)

	for _, d := range m.dispatchers {
		go m.startConsumer(d, m.shotdown)
	}

	e := <-m.shotdown
	return e
}

func (m *RabbitMQMessaging) newPubOpts(typ string) *PublishOpts {
	return &PublishOpts{
		Type:        typ,
		Count:       0,
		Traceparent: "without",
		MessageId:   uuid.New().String(),
		Delay:       time.Second,
	}
}

func (m *RabbitMQMessaging) installExchange(opt *ExchangeOpts) error {
	err := m.channel.ExchangeDeclare(opt.name, string(opt.kind), true, false, false, false, nil)

	if err != nil {
		return err
	}

	return nil
}

func (m *RabbitMQMessaging) installQueues(opts *QueueOpts) error {
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

func (m *RabbitMQMessaging) startConsumer(d *Dispatcher, shotdown chan error) {
	delivery, err := m.channel.Consume(d.Topology.Queue.Name, d.Topology.Binding.RoutingKey, false, false, false, false, nil)
	if err != nil {
		shotdown <- err
	}

	for received := range delivery {
		metadata, err := m.validateAndExtractMetadataFromDeliver(&received, d)
		if err != nil {
			received.Nack(true, false)
			continue
		}

		if metadata == nil {
			received.Nack(true, true)
			continue
		}

		m.logger.Info(LogMsgWithType("message received ", d.MsgType, received.MessageId))

		ptr := d.ReflectedType.Interface()
		err = json.Unmarshal(received.Body, ptr)
		if err != nil {
			m.logger.Error(LogMsgWithMessageId("unmarshal error", received.MessageId))
			received.Nack(true, false)
			continue
		}

		ctx, span, err := trace.SpanFromAMQPTraceparent(m.tracer, metadata.Traceparent, metadata.Type, d.Topology.Exchange.Name, d.Topology.Queue.Name)
		if err != nil {
			m.logger.Error("could not create a span")
			continue
		}

		if d.Topology.Queue.Retryable != nil && metadata.XCount > d.Topology.Queue.Retryable.NumberOfRetry {
			m.logger.Warn("message reprocessed to many times, sending to dead letter")
			received.Nack(true, false)
			span.SetStatus(codes.Error, "dead letter")
			span.End()
			continue
		}

		err = d.Handler(ptr, metadata)
		if err != nil {
			if d.Topology.Queue.Retryable == nil || err != errors.ErrorAMQPRetryable {
				received.Nack(true, false)
				span.SetStatus(codes.Error, "process failure without a retry")
				span.End()
				continue
			}

			m.logger.Warn(LogMessage("send message to process latter"))

			m.publishToDelayed(ctx, metadata, d.Topology, &received)

			received.Ack(true)
			span.End()
			continue
		}

		m.logger.Info(LogMsgWithMessageId("message processed properly", received.MessageId))
		received.Ack(true)
		span.SetStatus(codes.Ok, "success")
		span.End()
	}
}

func (m *RabbitMQMessaging) validateAndExtractMetadataFromDeliver(delivery *amqp.Delivery, d *Dispatcher) (*DeliveryMetadata, error) {
	msgID := delivery.MessageId
	if msgID == "" {
		m.logger.Error("unformatted amqp delivery - missing messageId parameter - send message to DLQ")
		return nil, errors.ErrorAMQPReceivedMessageValidator
	}

	typ := delivery.Type
	if typ == "" {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing type parameter - send message to DLQ", delivery.MessageId))
		return nil, errors.ErrorAMQPReceivedMessageValidator
	}

	xCount, ok := delivery.Headers[AMQPHeaderNumberOfRetry].(int64)
	if !ok {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing x-count header - send message to DLQ", delivery.MessageId))
		return nil, errors.ErrorAMQPReceivedMessageValidator
	}

	traceparent, ok := delivery.Headers[AMQPHeaderTraceparent]
	if !ok {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing x-trace-id header - send message to DLQ", delivery.MessageId))
		return nil, errors.ErrorAMQPReceivedMessageValidator
	}

	if typ != d.MsgType {
		return nil, nil
	}

	return &DeliveryMetadata{
		MessageId:   msgID,
		Type:        typ,
		XCount:      xCount,
		Traceparent: traceparent.(string),
		Headers:     delivery.Headers,
	}, nil
}

func (m *RabbitMQMessaging) publishToDelayed(ctx context.Context, metadata *DeliveryMetadata, t *Topology, received *amqp.Delivery) error {

	return m.ch.Publish(t.delayed.ExchangeName, t.delayed.RoutingKey, false, false, amqp.Publishing{
		Headers: amqp.Table{
			AMQPHeaderNumberOfRetry: metadata.XCount + 1,
			AMQPHeaderTraceparent:   metadata.Traceparent,
			AMQPHeaderDelay:         t.Queue.Retryable.DelayBetween.Milliseconds(),
		},
		Type:        received.Type,
		ContentType: received.ContentType,
		MessageId:   received.MessageId,
		UserId:      received.UserId,
		AppId:       received.AppId,
		Body:        received.Body,
	})
}
