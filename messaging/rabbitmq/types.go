package rabbitmq

import (
	"context"
	"reflect"

	"github.com/streadway/amqp"

	"github.com/ralvescostati/pkgs/logger"
)

type (
	ExchangeKind string

	// Params is a RabbitMQ params needed to declare an Exchange, Queue or Bind them
	Params struct {
		ExchangeName     string
		ExchangeType     ExchangeKind
		QueueName        string
		RoutingKey       string
		Retryable        bool
		EnabledTelemetry bool
	}

	PublishOpts struct {
		Key   string
		Value any
	}

	SubHandler = func(msg any, opts map[string]any) error

	// IRabbitMQMessaging is RabbitMQ Config Builder
	IRabbitMQMessaging interface {
		// AssertExchange Declare a durable, not excluded Exchange with the following parameters
		AssertExchange(params *Params) IRabbitMQMessaging

		// AssertExchangeAssertQueue Declare a durable, not excluded Queue with the following parameters
		AssertQueue(params *Params) IRabbitMQMessaging

		// Binding bind an exchange/queue with the following parameters without extra RabbitMQ configurations such as Dead Letter.
		Binding(params *Params) IRabbitMQMessaging

		// AssertExchange Declare a durable, not excluded Exchange with the following parameters with a default Dead Letter exchange
		AssertExchangeWithDeadLetter() IRabbitMQMessaging

		// AssertDelayedExchange will be declare a Delay exchange and configure a dead letter exchange and queue.
		//
		// When messages for delay exchange was noAck these messages will sent to the dead letter exchange/queue.
		AssertDelayedExchange() IRabbitMQMessaging

		Publisher(ctx context.Context, params *Params, msg any, opts ...PublishOpts) error
		Subscriber(ctx context.Context, params *Params) error

		// AddDispatcher Add the handler and msg type
		//
		// Each time a message came, we check the queue, and get the available handlers for that queue.
		// After we do a coercion of the msg type to check which handler expect this msg type
		AddDispatcher(event string, handler SubHandler, structWillUseToTypeCoercion any) error

		Build() (IRabbitMQMessaging, error)
	}

	// AMQPChannel is an abstraction for AMQP default channel to improve unit tests
	AMQPChannel interface {
		ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
		QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
		QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
		Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	}

	Dispatcher struct {
		Queue          string
		ReceiveMsgType string
		ReflectedType  reflect.Value
		Handler        SubHandler
	}

	// IRabbitMQMessaging is the implementation for IRabbitMQMessaging
	RabbitMQMessaging struct {
		Err         error
		logger      logger.ILogger
		conn        *amqp.Connection
		ch          AMQPChannel
		dispatchers map[string][]*Dispatcher
	}
)
