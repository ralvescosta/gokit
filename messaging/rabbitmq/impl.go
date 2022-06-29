package rabbitmq

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/streadway/amqp"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logger"
)

// New(...) create a new instance for IRabbitMQMessaging
//
// New(...) connect to the RabbitMQ broker and stablish a channel
func New(cfg *env.Configs, logger logger.ILogger) IRabbitMQMessaging {
	rb := &RabbitMQMessaging{
		logger:      logger,
		config:      cfg,
		dispatchers: []*Dispatcher{},
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.RABBIT_USER, cfg.RABBIT_PASSWORD, cfg.RABBIT_VHOST, cfg.RABBIT_PORT))
	if err != nil {
		logger.Error(fmt.Sprintf(ConnErrorMessage, "broker", err))
		rb.Err = err
		return rb
	}

	rb.conn = conn
	ch, err := conn.Channel()
	if err != nil {
		logger.Error(fmt.Sprintf(ConnErrorMessage, "channel", err))
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

	m.queuesToDeclare = append(m.queuesToDeclare, params)

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

func (m *RabbitMQMessaging) Build() (IRabbitMQMessaging, error) {
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

func (m *RabbitMQMessaging) RegisterDispatcher(queue string, handler SubHandler, structWillUseToTypeCoercion any) error {
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
	// delivery, err := m.ch.Consume(params.QueueName, params.RoutingKey, false, false, false, false, nil)
	// if err != nil {
	// 	return err
	// }

	shotdown := make(chan error)

	// go m.exec(params, delivery)
	for _, d := range m.dispatchers {
		go m.startConsumer(d, shotdown)
	}

	select {
	case e := <-shotdown:
		return e
	}

	return nil
}

func (m *RabbitMQMessaging) startConsumer(d *Dispatcher, shotdown chan error) {
	delivery, err := m.ch.Consume(d.DeclareParams.QueueName, d.BindParams.RoutingKey, false, false, false, false, nil)
	if err != nil {
		shotdown <- err
	}

	for received := range delivery {
		shouldSkip, err := m.validateDelivered(&received, d)
		//remove msg from queue and send to DLQ
		if err != nil {
			received.Nack(true, true)
		}

		//skip msg
		if shouldSkip {
			received.Ack(true)
		}

		// ptr := d.ReflectedType.Interface()
	}
}

func (m *RabbitMQMessaging) validateDelivered(delivery *amqp.Delivery, d *Dispatcher) (bool, error) {
	typ := delivery.Type
	if typ == "" {
		delivery.Headers["rejection-reason"] = "RequiredParameter:Type"
		return false, errors.New("")
	}

	// xCount, ok := delivery.Headers["x-count"]
	// if d.DeclareParams.Retryable != nil && !ok {
	// 	//nack
	// }

	return false, nil
}

// func (m *RabbitMQMessaging) exec(params *Params, delivery <-chan amqp.Delivery) {
// 	for received := range delivery {
// 		if received.Type == "" {
// 			m.logger.Warn("[RabbitMQ:HandlerExecutor] ignore message reason: message without type header")
// 			received.Ack(true)
// 			continue
// 		}

// 		dispatchers, ok := m.dispatchers[params.QueueName]
// 		if !ok {
// 			m.logger.Warn("[RabbitMQ:HandlerExecutor] ignore message reason: there is no handler for this queue registered yet")
// 			received.Ack(true)
// 			continue
// 		}

// 		var mPointer any
// 		var handler SubHandler

// 		for _, d := range dispatchers {
// 			if d.ReceiveMsgType == received.Type {
// 				mPointer = d.ReflectedType.Interface()

// 				err := json.Unmarshal(received.Body, mPointer)
// 				if err == nil {
// 					handler = d.Handler
// 					break
// 				}
// 			}
// 		}

// 		if mPointer == nil || handler == nil {
// 			m.logger.Error(fmt.Sprintf("[RabbitMQ:HandlerExecutor] ignore message reason: failure type coercion. Queue: %s.", params.QueueName))
// 			received.Ack(true)
// 			continue
// 		}

// 		m.logger.Info(fmt.Sprintf("[RabbitMQ:HandlerExecutor] message received %T", mPointer))

// 		err := handler(mPointer, nil)
// 		if err == nil {
// 			m.logger.Info("[RabbitMQ:HandlerExecutor] message properly processed")
// 			received.Ack(true)
// 			continue
// 		}

// 		m.logger.Error(err.Error())

// 		if params.Retryable != nil {
// 			m.logger.Warn("[RabbitMQ:HandlerExecutor] message has no retry police, purging message")
// 			received.Ack(true)
// 			continue
// 		}

// 		m.logger.Debug("[RabbitMQ:HandlerExecutor] sending failure msg to delayed exchange")
// 		m.Publisher(context.Background(), nil, nil)

// 		received.Ack(true)
// 	}
// }
