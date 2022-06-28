package rabbitmq

import (
	"context"
	"reflect"
	"time"

	"github.com/streadway/amqp"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logger"
)

type (
	ExchangeKind string

	// Params is a RabbitMQ params needed to declare an Exchange, Queue or Bind them
	Params struct {
		ExchangeName   string
		ExchangeType   ExchangeKind
		QueueName      string
		QueueTTL       time.Duration
		RoutingKey     string
		Retryable      *Retry
		WithDeadLatter bool
	}

	Retry struct {
		TTL time.Duration
	}

	PublishOpts struct {
		Key   string
		Value any
	}

	SubHandler = func(msg any, opts map[string]any) error

	// IRabbitMQMessaging is RabbitMQ Config Builder
	IRabbitMQMessaging interface {
		// AssertExchange Declare a durable, not excluded Exchange with the following parameters
		DeclareExchange(params *Params) IRabbitMQMessaging

		// AssertExchangeAssertQueue Declare a durable, not excluded Queue with the following parameters
		DeclareQueue(params *Params) IRabbitMQMessaging

		// Binding bind an exchange/queue with the following parameters without extra RabbitMQ configurations such as Dead Letter.
		BindQueue(params *Params) IRabbitMQMessaging

		// AssertQueueWithDeadLetter Declare a durable, not excluded Exchange with the following parameters with a default Dead Letter exchange
		// DeclareQueueWithDeadLetter(params *Params) IRabbitMQMessaging

		// AssertDelayedExchange will be declare a Delay exchange and configure a dead letter exchange and queue.
		//
		// When messages for delay exchange was noAck these messages will sent to the dead letter exchange/queue.
		// DeclareDelayedExchange(params *Params) IRabbitMQMessaging

		Publisher(ctx context.Context, params *Params, msg any, opts ...PublishOpts) error
		Consumer() error

		// AddDispatcher Add the handler and msg type
		//
		// Each time a message came, we check the queue, and get the available handlers for that queue.
		// After we do a coercion of the msg type to check which handler expect this msg type
		RegisterDispatcher(event string, handler SubHandler, structWillUseToTypeCoercion any) error

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
		Queue          string
		ReceiveMsgType string
		ReflectedType  reflect.Value
		Handler        SubHandler
	}

	// IRabbitMQMessaging is the implementation for IRabbitMQMessaging
	RabbitMQMessaging struct {
		Err                error
		logger             logger.ILogger
		config             *env.Configs
		conn               *amqp.Connection
		ch                 AMQPChannel
		exchangeToDeclare  []*Params
		queueToDeclare     []*Params
		exchangesToBinding []*Params
		queueToBinding     []*Params
		ttl                time.Duration
		dispatchers        map[string][]*Dispatcher
	}
)
