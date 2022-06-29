package rabbitmq

import (
	"reflect"
	"time"

	"github.com/streadway/amqp"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logger"
)

type (
	ExchangeKind string

	Retry struct {
		NumberOfRetry int
		DelayBetween  time.Duration
	}

	// Params is a RabbitMQ params needed to declare an Exchange, Queue or Bind them
	DeclareExchangeParams struct {
		ExchangeName string
		ExchangeType ExchangeKind
	}

	DeclareQueueParams struct {
		QueueName      string
		QueueTTL       time.Duration
		Retryable      *Retry
		WithDeadLatter bool
	}

	BindExchangeParams struct {
		ExchangeName string
		RoutingKey   string
	}

	BindQueueParams struct {
		QueueName    string
		ExchangeName string
		RoutingKey   string
	}

	PublishOpts struct {
		Key   string
		Value any
	}

	DeliveryMetadata struct {
		MessageId string
		XCount    int
		Type      string
		TraceId   string
		Headers   map[string]interface{}
	}

	ConsumerHandler = func(msg any, metadata *DeliveryMetadata) error

	// IRabbitMQMessaging is RabbitMQ Config Builder
	IRabbitMQMessaging interface {
		// AssertExchange Declare a durable, not excluded Exchange with the following parameters
		DeclareExchange(params *DeclareExchangeParams) IRabbitMQMessaging

		// AssertExchangeAssertQueue Declare a durable, not excluded Queue with the following parameters
		DeclareQueue(params *DeclareQueueParams) IRabbitMQMessaging

		// Binding bind an exchange/queue with the following parameters without extra RabbitMQ configurations such as Dead Letter.
		BindQueue(params *BindQueueParams) IRabbitMQMessaging

		// AssertQueueWithDeadLetter Declare a durable, not excluded Exchange with the following parameters with a default Dead Letter exchange
		// DeclareQueueWithDeadLetter(params *Params) IRabbitMQMessaging

		// AssertDelayedExchange will be declare a Delay exchange and configure a dead letter exchange and queue.
		//
		// When messages for delay exchange was noAck these messages will sent to the dead letter exchange/queue.
		// DeclareDelayedExchange(params *Params) IRabbitMQMessaging

		Publisher() error
		Consume() error

		// AddDispatcher Add the handler and msg type
		//
		// Each time a message came, we check the queue, and get the available handlers for that queue.
		// After we do a coercion of the msg type to check which handler expect this msg type
		RegisterDispatcher(event string, handler ConsumerHandler, structWillUseToTypeCoercion any) error

		Build() (IRabbitMQMessaging, error)
	}

	// AMQPChannel is an abstraction for AMQP default channel to improve unit tests
	AMQPChannel interface {
		ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
		QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
		QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
		Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
		Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	}

	Dispatcher struct {
		Queue         string
		BindParams    *BindQueueParams
		DeclareParams *DeclareQueueParams
		MsgType       string
		ReflectedType reflect.Value
		Handler       ConsumerHandler
	}

	// IRabbitMQMessaging is the implementation for IRabbitMQMessaging
	RabbitMQMessaging struct {
		Err                error
		logger             logger.ILogger
		config             *env.Configs
		conn               *amqp.Connection
		ch                 AMQPChannel
		exchangesToDeclare []*DeclareExchangeParams
		queuesToDeclare    []*DeclareQueueParams
		exchangesToBinding []*BindExchangeParams
		queuesToBinding    []*BindQueueParams
		dispatchers        []*Dispatcher
	}
)
