package rabbitmq

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logging"
)

// New(...) create a new instance for IRabbitMQMessaging
//
// New(...) connect to the RabbitMQ broker and stablish a channel
func New(cfg *env.Configs, logger logging.ILogger) IRabbitMQMessaging {
	rb := &RabbitMQMessaging{
		logger:      logger,
		config:      cfg,
		dispatchers: []*Dispatcher{},
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.RABBIT_USER, cfg.RABBIT_PASSWORD, cfg.RABBIT_VHOST, cfg.RABBIT_PORT))
	if err != nil {
		logger.Error(LogMessage("failure to connect to the broker"), logging.ErrorField(err))
		rb.Err = ErrorConnection
		return rb
	}

	rb.conn = conn
	ch, err := conn.Channel()
	if err != nil {
		logger.Error(LogMessage("failure to establish the channel"), logging.ErrorField(err))
		rb.Err = ErrorChannel
		return rb
	}

	rb.ch = ch

	return rb
}

func (m *RabbitMQMessaging) DeclareExchange(params *DeclareExchangeParams) IRabbitMQMessaging {
	if m.Err != nil {
		return m
	}

	m.exchangesToDeclare = append(m.exchangesToDeclare, params)

	return m
}

func (m *RabbitMQMessaging) DeclareQueue(params *DeclareQueueParams) IRabbitMQMessaging {
	if m.Err != nil {
		return m
	}

	m.queuesToDeclare = append(m.queuesToDeclare, params)
	if params.WithDeadLatter || params.Retryable != nil {
		dlqQueue := m.newFallbackQueueName(DLQ_FALLBACK, params.QueueName)
		dlqExchange := m.newFallbackExchangeName(DLQ_FALLBACK, params.QueueName)

		m.queuesToDeclare = append(m.queuesToDeclare, &DeclareQueueParams{
			QueueName: dlqQueue,
		})

		m.exchangesToDeclare = append(m.exchangesToDeclare, &DeclareExchangeParams{
			ExchangeName: dlqExchange,
			ExchangeType: DIRECT_EXCHANGE,
		})

		m.queuesToBinding = append(m.queuesToBinding, &BindQueueParams{
			QueueName:    dlqQueue,
			ExchangeName: dlqExchange,
			RoutingKey:   m.newFallbackRoutingKey(DLQ_FALLBACK, params.QueueName),
		})
	}

	if params.Retryable != nil {
		delayedExchange := m.newFallbackExchangeName(RETRY_FALLBACK, params.QueueName)

		m.exchangesToDeclare = append(m.exchangesToDeclare, &DeclareExchangeParams{
			ExchangeName: delayedExchange,
			ExchangeType: DELAY_EXCHANGE,
		})

		m.queuesToBinding = append(m.queuesToBinding, &BindQueueParams{
			QueueName:    params.QueueName,
			ExchangeName: delayedExchange,
			RoutingKey:   m.newFallbackRoutingKey(RETRY_FALLBACK, params.QueueName),
		})
	}

	return m
}

func (m *RabbitMQMessaging) BindExchange(params *BindExchangeParams) IRabbitMQMessaging {
	if m.Err != nil {
		return m
	}

	m.exchangesToBinding = append(m.exchangesToBinding, params)

	return m
}

func (m *RabbitMQMessaging) BindQueue(params *BindQueueParams) IRabbitMQMessaging {
	if m.Err != nil {
		return m
	}

	m.queuesToBinding = append(m.queuesToBinding, params)

	return m
}

func (m *RabbitMQMessaging) Build() (IRabbitMQMessaging, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	m.logger.Debug(LogMessage("declaring exchanges..."))
	for _, exch := range m.exchangesToDeclare {
		if err := m.declareExchange(exch); err != nil {
			m.logger.Error(LogMessage("declare exchange err"), logging.ErrorField(err))
			return nil, err
		}
	}
	m.logger.Debug(LogMessage("exchanges declared"))

	m.logger.Debug(LogMessage("binding exchanges..."))
	for _, exch := range m.exchangesToBinding {
		if err := m.bindExchanges(exch); err != nil {
			m.logger.Error(LogMessage("bind exchange err"), logging.ErrorField(err))
			return nil, err
		}
	}
	m.logger.Debug(LogMessage("exchanges bound"))

	m.logger.Debug(LogMessage("declaring queues..."))
	for _, q := range m.queuesToDeclare {
		if err := m.declareQueue(q); err != nil {
			m.logger.Error(LogMessage("declare queue err"), logging.ErrorField(err))
			return nil, err
		}
	}
	m.logger.Debug(LogMessage("queues declared"))

	m.logger.Debug(LogMessage("binding queues..."))
	for _, q := range m.queuesToBinding {
		if err := m.bindQueue(q); err != nil {
			m.logger.Error(LogMessage("bind queue err"), logging.ErrorField(err))
			return nil, err
		}
	}
	m.logger.Debug(LogMessage("queues bound"))

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
			AMQPHeaderTraceID:       opts.TraceId,
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

func (m *RabbitMQMessaging) RegisterDispatcher(queue string, handler ConsumerHandler, structWillUseToTypeCoercion any) error {
	if structWillUseToTypeCoercion == nil || queue == "" {
		return ErrorRegisterDispatcher
	}

	var bind *BindQueueParams
	var declare *DeclareQueueParams

	for _, v := range m.queuesToBinding {
		if v.QueueName == queue {
			bind = v
			break
		}
	}

	for _, v := range m.queuesToDeclare {
		if v.QueueName == queue {
			declare = v
		}
	}

	dispatch := &Dispatcher{
		Queue:         queue,
		BindParams:    bind,
		DeclareParams: declare,
		Handler:       handler,
		MsgType:       fmt.Sprintf("%T", structWillUseToTypeCoercion),
		ReflectedType: reflect.New(reflect.TypeOf(structWillUseToTypeCoercion).Elem()),
	}

	m.dispatchers = append(m.dispatchers, dispatch)

	return nil
}

func (m *RabbitMQMessaging) Consume() error {
	if m.Err != nil {
		return m.Err
	}

	shotdown := make(chan error)

	for _, d := range m.dispatchers {
		go m.startConsumer(d, shotdown)
	}

	e := <-shotdown
	return e
}

func (m *RabbitMQMessaging) newPubOpts(typ string) *PublishOpts {
	return &PublishOpts{
		Type:      typ,
		Count:     0,
		TraceId:   "without",
		MessageId: uuid.NewString(),
		Delay:     time.Second,
	}
}

func (m *RabbitMQMessaging) newRoutingKey(exchange, queue string) string {
	return exchange + "-" + queue
}

func (m *RabbitMQMessaging) newFallbackExchangeName(typ FallbackType, queue string) string {
	return string(typ) + "-" + queue
}

func (m *RabbitMQMessaging) newFallbackQueueName(typ FallbackType, queue string) string {
	return string(typ) + "-" + queue
}

func (m *RabbitMQMessaging) newFallbackRoutingKey(typ FallbackType, queue string) string {
	return string(typ) + "-" + queue + "-key"
}

func (m *RabbitMQMessaging) declareExchange(params *DeclareExchangeParams) error {
	var args amqp.Table
	if params.ExchangeType == DELAY_EXCHANGE {
		args = amqp.Table{
			"x-delayed-type": "direct",
		}
	}

	return m.ch.ExchangeDeclare(params.ExchangeName, string(params.ExchangeType), true, false, false, false, args)
}

func (m *RabbitMQMessaging) bindExchanges(params *BindExchangeParams) error {
	for _, e := range params.ExchangesDestinations {
		err := m.ch.ExchangeBind(e, m.newRoutingKey(params.ExchangeSource, e), params.ExchangeSource, false, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *RabbitMQMessaging) declareQueue(params *DeclareQueueParams) error {
	var amqpTable amqp.Table
	if params.WithDeadLatter || params.Retryable != nil {
		amqpTable = amqp.Table{
			"x-dead-letter-exchange":    m.newFallbackExchangeName(DLQ_FALLBACK, params.QueueName),
			"x-dead-letter-routing-key": m.newFallbackRoutingKey(DLQ_FALLBACK, params.QueueName),
		}
	}

	_, err := m.ch.QueueDeclare(params.QueueName, true, false, false, false, amqpTable)
	if err != nil {
		return err
	}

	return nil
}

func (m *RabbitMQMessaging) bindQueue(params *BindQueueParams) error {
	routingKey := params.RoutingKey
	if routingKey == "" {
		routingKey = m.newRoutingKey(params.ExchangeName, params.QueueName)
	}

	err := m.ch.QueueBind(params.QueueName, routingKey, params.ExchangeName, false, nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *RabbitMQMessaging) startConsumer(d *Dispatcher, shotdown chan error) {
	delivery, err := m.ch.Consume(d.DeclareParams.QueueName, m.newRoutingKey(d.BindParams.QueueName, d.BindParams.ExchangeName), false, false, false, false, nil)
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
			m.logger.Debug(LogMsgWithMessageId("skipping amqp delivery - different msg type - send back to queue", received.MessageId))
			received.Nack(true, true)
			continue
		}

		ptr := d.ReflectedType.Interface()
		err = json.Unmarshal(received.Body, ptr)
		if err != nil {
			m.logger.Error(LogMsgWithMessageId("unmarshal error", received.MessageId))
			received.Nack(true, false)
			continue
		}

		if d.DeclareParams.Retryable != nil && metadata.XCount > d.DeclareParams.Retryable.NumberOfRetry {
			m.logger.Warn("message reprocessed to many times, sending to dead letter")
			received.Nack(true, false)
			continue
		}

		m.logger.Info(LogMsgWithType("message received ", d.MsgType, received.MessageId))

		err = d.Handler(ptr, metadata)
		if err != nil {
			if d.DeclareParams.Retryable == nil || err != ErrorRetryable {
				received.Nack(true, false)
				continue
			}

			m.logger.Warn(LogMessage("send message to process latter"))

			m.publishToDelayed(metadata, d.BindParams, d.DeclareParams, &received)

			received.Ack(true)
			continue
		}

		m.logger.Info(LogMsgWithMessageId("message processed properly", received.MessageId))
		received.Ack(true)
	}
}

func (m *RabbitMQMessaging) validateAndExtractMetadataFromDeliver(delivery *amqp.Delivery, d *Dispatcher) (*DeliveryMetadata, error) {
	msgID := delivery.MessageId
	if msgID != "" {
		m.logger.Error("unformatted amqp delivery - missing messageId parameter - send message to DLQ")
		return nil, ErrorReceivedMessageValidator
	}

	typ := delivery.Type
	if typ == "" {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing type parameter - send message to DLQ", delivery.MessageId))
		return nil, ErrorReceivedMessageValidator
	}

	xCount, ok := delivery.Headers[AMQPHeaderNumberOfRetry].(int64)
	if !ok {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing x-count header - send message to DLQ", delivery.MessageId))
		return nil, ErrorReceivedMessageValidator
	}

	traceID, ok := delivery.Headers[AMQPHeaderTraceID]
	if !ok {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing x-trace-id header - send message to DLQ", delivery.MessageId))
		return nil, ErrorReceivedMessageValidator
	}

	if typ != d.MsgType {
		return nil, nil
	}

	return &DeliveryMetadata{
		MessageId: msgID,
		Type:      typ,
		XCount:    xCount,
		TraceId:   traceID.(string),
		Headers:   delivery.Headers,
	}, nil
}

func (m *RabbitMQMessaging) publishToDelayed(metadata *DeliveryMetadata, b *BindQueueParams, d *DeclareQueueParams, received *amqp.Delivery) error {
	exch := m.newFallbackExchangeName(RETRY_FALLBACK, b.QueueName)
	rk := m.newFallbackRoutingKey(RETRY_FALLBACK, b.QueueName)

	return m.ch.Publish(exch, rk, false, false, amqp.Publishing{
		Headers: amqp.Table{
			AMQPHeaderNumberOfRetry: metadata.XCount + 1,
			AMQPHeaderTraceID:       metadata.TraceId,
			AMQPHeaderDelay:         d.Retryable.DelayBetween.Milliseconds(),
		},
		Type:        received.Type,
		ContentType: received.ContentType,
		MessageId:   received.MessageId,
		UserId:      received.UserId,
		AppId:       received.AppId,
		Body:        received.Body,
	})
}
