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
		topologies:  []*Topology{},
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

	rb.ch = ch

	return rb
}

var dial = func(cfg *env.Configs) (AMQPConnection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.RABBIT_USER, cfg.RABBIT_PASSWORD, cfg.RABBIT_VHOST, cfg.RABBIT_PORT))
}

func (m *RabbitMQMessaging) Declare(opts *Topology) IRabbitMQMessaging {
	if m.Err != nil {
		return m
	}

	if opts.isBindable {
		m.bind(opts)
	}

	m.topologies = append(m.topologies, opts)

	return m
}

func (m *RabbitMQMessaging) ApplyBinds() IRabbitMQMessaging {
	if m.Err != nil {
		return m
	}

	for _, v := range m.topologies {
		m.bind(v)
	}

	return m
}

func (m *RabbitMQMessaging) bind(params *Topology) {
	params.Binding = m.newBinding(params)
	params.deadLetter = m.newDeadLetter(params)
	params.delayed = m.newDelayed(params)

	if params.deadLetter != nil {
		params.Binding.dlqRoutingKey = params.deadLetter.RoutingKey
	}

	if params.delayed != nil {
		params.Binding.delayedRoutingKey = params.delayed.RoutingKey
	}
}

func (m *RabbitMQMessaging) newBinding(params *Topology) *BindingOpts {
	return &BindingOpts{
		RoutingKey: m.newRoutingKey(params.Exchange.Name, params.Queue.Name),
	}
}

func (m *RabbitMQMessaging) newDeadLetter(params *Topology) *DeadLetterOpts {
	if !params.Queue.WithDeadLatter && params.Queue.Retryable == nil {
		return nil
	}

	return &DeadLetterOpts{
		QueueName:    m.newFallbackName(DLQ_FALLBACK, params.Queue.Name),
		ExchangeName: params.Exchange.Name,
		RoutingKey:   m.newFallbackName(DLQ_FALLBACK, params.Binding.RoutingKey),
	}
}

func (m *RabbitMQMessaging) newDelayed(params *Topology) *DelayedOpts {
	if params.Queue.Retryable == nil {
		return nil
	}

	return &DelayedOpts{
		QueueName:    params.Queue.Name,
		ExchangeName: m.newFallbackName(RETRY_FALLBACK, params.Exchange.Name),
		RoutingKey:   m.newFallbackName(RETRY_FALLBACK, params.Binding.RoutingKey),
	}
}

func (m *RabbitMQMessaging) Build() (IRabbitMQMessaging, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	for _, d := range m.topologies {
		m.logger.Debug(LogMessage("declaring exchanges..."))
		if err := m.declareExchange(d); err != nil {
			m.logger.Error(LogMessage("declare exchange err"), logging.ErrorField(err))
			return nil, err
		}
		m.logger.Debug(LogMessage("exchanges declared"))

		m.logger.Debug(LogMessage("binding exchanges..."))
		if err := m.bindExchanges(d); err != nil {
			m.logger.Error(LogMessage("bind exchange err"), logging.ErrorField(err))
			return nil, err
		}
		m.logger.Debug(LogMessage("exchanges bound"))

		m.logger.Debug(LogMessage("declaring queues..."))
		if err := m.declareQueue(d); err != nil {
			m.logger.Error(LogMessage("declare queue err"), logging.ErrorField(err))
			return nil, err
		}
		m.logger.Debug(LogMessage("queues declared"))

		m.logger.Debug(LogMessage("binding queues..."))
		if err := m.bindQueue(d); err != nil {
			m.logger.Error(LogMessage("bind queue err"), logging.ErrorField(err))
			return nil, err
		}
		m.logger.Debug(LogMessage("queues bound"))
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

	return m.ch.Publish(exchange, routingKey, false, false, amqp.Publishing{
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

func (m *RabbitMQMessaging) newRoutingKey(exchange, queue string) string {
	return exchange + "-" + queue + "-" + "key"
}

func (m *RabbitMQMessaging) newFallbackName(typ FallbackType, name string) string {
	return string(typ) + "-" + name
}

func (m *RabbitMQMessaging) declareExchange(opt *Topology) error {
	if opt.Exchange != nil {
		err := m.ch.ExchangeDeclare(opt.Exchange.Name, string(opt.Exchange.Type), true, false, false, false, nil)
		if err != nil {
			return err
		}
	}

	if opt.delayed == nil {
		return nil
	}

	err := m.ch.ExchangeDeclare(opt.delayed.ExchangeName, string(DELAY_EXCHANGE), true, false, false, false, amqp.Table{
		"x-delayed-type": "direct",
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *RabbitMQMessaging) bindExchanges(opts *Topology) error {
	if opts.Exchange.Bindings == nil || len(opts.Exchange.Bindings) == 0 {
		return nil
	}

	for _, e := range opts.Exchange.Bindings {
		err := m.ch.ExchangeBind(e, m.newRoutingKey(opts.Exchange.Name, e), opts.Exchange.Name, false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *RabbitMQMessaging) declareQueue(opts *Topology) error {
	if opts.Queue == nil {
		return nil
	}

	var amqpTable amqp.Table
	if opts.deadLetter != nil || opts.delayed != nil {
		//when we do not specify the exchange and configure in the dlq routing the queue name
		//when messages was rejected will be sent to dql queue directly
		amqpTable = amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": opts.deadLetter.QueueName,
		}

		_, err := m.ch.QueueDeclare(opts.deadLetter.QueueName, true, false, false, false, nil)
		if err != nil {
			return err
		}
	}

	_, err := m.ch.QueueDeclare(opts.Queue.Name, true, false, false, false, amqpTable)
	if err != nil {
		return err
	}

	return nil
}

func (m *RabbitMQMessaging) bindQueue(opts *Topology) error {
	if err := m.ch.QueueBind(opts.Queue.Name, opts.Binding.RoutingKey, opts.Exchange.Name, false, nil); err != nil {
		return err
	}

	if opts.delayed != nil {
		if err := m.ch.QueueBind(opts.delayed.QueueName, opts.Binding.delayedRoutingKey, opts.delayed.ExchangeName, false, nil); err != nil {
			return err
		}
	}

	return nil
}

func (m *RabbitMQMessaging) startConsumer(d *Dispatcher, shotdown chan error) {
	delivery, err := m.ch.Consume(d.Topology.Queue.Name, d.Topology.Binding.RoutingKey, false, false, false, false, nil)
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
