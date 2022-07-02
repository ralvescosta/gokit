package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

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
		rb.Err = err
		return rb
	}

	rb.conn = conn
	ch, err := conn.Channel()
	if err != nil {
		logger.Error(LogMessage("failure to establish the channel"), logging.ErrorField(err))
		rb.Err = err
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

	//if there are DLQ create a new declare and binding to dlq
	//if there are retry same thing
	m.queuesToDeclare = append(m.queuesToDeclare, params)
	if params.WithDeadLatter || params.Retryable != nil {
		m.queuesToDeclare = append(m.queuesToDeclare, &DeclareQueueParams{
			QueueName: m.newFallbackQueueName(DLQ_FALLBACK, params.QueueName),
		})
		m.exchangesToDeclare = append(m.exchangesToDeclare, &DeclareExchangeParams{
			ExchangeName: m.newFallbackExchangeName(DLQ_FALLBACK, params.QueueName),
			ExchangeType: DIRECT_EXCHANGE,
		})
		m.queuesToBinding = append(m.queuesToBinding, &BindQueueParams{
			QueueName:    m.newFallbackQueueName(DLQ_FALLBACK, params.QueueName),
			ExchangeName: m.newFallbackExchangeName(DLQ_FALLBACK, params.QueueName),
			RoutingKey:   m.newFallbackRoutingKey(DLQ_FALLBACK, params.QueueName),
		})
	}

	if params.Retryable != nil {
		m.exchangesToDeclare = append(m.exchangesToDeclare, &DeclareExchangeParams{
			ExchangeName: m.newFallbackExchangeName(RETRY_FALLBACK, params.QueueName),
			ExchangeType: DELAY_EXCHANGE,
		})
		m.queuesToBinding = append(m.queuesToBinding, &BindQueueParams{
			QueueName:    params.QueueName,
			ExchangeName: m.newFallbackExchangeName(RETRY_FALLBACK, params.QueueName),
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
	m.logger.Debug(LogMessage("exchanges binded"))

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
	m.logger.Debug(LogMessage("queues binded"))

	return m, m.Err
}

func (m *RabbitMQMessaging) Publisher() error {
	// byt, err := json.Marshal(msg)
	// if err != nil {
	// 	m.logger.Error(err.Error())
	// 	return err
	// }

	// err = m.ch.Publish(params.ExchangeName, params.RoutingKey, true, true, amqp.Publishing{
	// 	AppId:       m.config.APP_NAME,
	// 	MessageId:   uuid.NewString(),
	// 	ContentType: JsonContentType,
	// 	Type:        fmt.Sprintf("%T", msg),
	// 	Timestamp:   time.Now(),
	// 	UserId:      m.config.RABBIT_USER,
	// 	// Headers: amqp.Table{},
	// 	Body: byt,
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (m *RabbitMQMessaging) RegisterDispatcher(queue string, handler ConsumerHandler, structWillUseToTypeCoercion any) error {
	if structWillUseToTypeCoercion == nil || queue == "" {
		return errors.New("[RabbitMQ:AddDispatcher]")
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
	shotdown := make(chan error)

	for _, d := range m.dispatchers {
		go m.startConsumer(d, shotdown)
	}

	e := <-shotdown
	return e
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

		m.logger.Info(LogMsgWithType("message received ", d.MsgType, received.MessageId))

		err = d.Handler(ptr, metadata)
		if err != nil {
			//retry flow - ack the msg, encrise the xCount, publish to delayed exchange
			//no retry flow - nack the msg
			if d.DeclareParams.Retryable != nil {
				if metadata.XCount < d.DeclareParams.Retryable.NumberOfRetry {
					m.ch.Publish(m.newFallbackExchangeName(RETRY_FALLBACK, d.BindParams.QueueName), m.newFallbackRoutingKey(RETRY_FALLBACK, d.BindParams.QueueName), false, false, amqp.Publishing{
						Headers: amqp.Table{
							AMQPHeaderNumberOfRetry: metadata.XCount + 1,
							AMQPHeaderTraceID:       metadata.TraceId,
							AMQPHeaderDelay:         d.DeclareParams.Retryable.DelayBetween,
						},
						ContentType: received.ContentType,
						MessageId:   received.MessageId,
						UserId:      received.UserId,
						AppId:       received.AppId,
						Body:        received.Body,
					})
					received.Ack(true)
				} else {
					received.Nack(true, false)
				}
			}
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
		return nil, errors.New("")
	}

	typ := delivery.Type
	if typ == "" {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing type parameter - send message to DLQ", delivery.MessageId))
		return nil, errors.New("")
	}

	xCountHeader, ok := delivery.Headers[AMQPHeaderNumberOfRetry]
	xCount, err := strconv.Atoi(xCountHeader.(string))
	if !ok || err != nil {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing x-count header - send message to DLQ", delivery.MessageId))
		return nil, errors.New("")
	}

	traceID, ok := delivery.Headers[AMQPHeaderTraceID]
	if !ok {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing x-trace-id header - send message to DLQ", delivery.MessageId))
		return nil, errors.New("")
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

// func (m *RabbitMQMessaging) DeclareQueueWithDeadLetter(params *Params) IRabbitMQMessaging {
// 	if m.Err != nil {
// 		return m
// 	}

// 	if params.ExchangeName == "" || params.RoutingKey == "" {
// 		m.Err = errors.New("")
// 		return m
// 	}

// 	_, err := m.ch.QueueDeclare(params.QueueName, true, false, false, false, amqp.Table{
// 		"x-dead-letter-exchange":    fmt.Sprintf("%s%s", params.ExchangeName, DeadLetterSuffix),
// 		"x-dead-letter-routing-key": fmt.Sprintf("%s%s", params.RoutingKey, DeadLetterSuffix),
// 	})
// 	if err != nil {
// 		m.Err = err
// 		return m
// 	}

// 	return m
// }
