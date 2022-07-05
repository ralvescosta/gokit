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
	FallbackType string

	// Retry
	Retry struct {
		NumberOfRetry int64
		DelayBetween  time.Duration
	}

	// QueueOpts declare queue configuration
	QueueOpts struct {
		Name           string
		TTL            time.Duration
		Retryable      *Retry
		WithDeadLatter bool
	}

	// ExchangeOpts exchanges to declare
	ExchangeOpts struct {
		Name     string
		Type     ExchangeKind
		Bindings []string
	}

	// BindingOpts binds configuration
	BindingOpts struct {
		RoutingKey        string
		dlqRoutingKey     string
		delayedRoutingKey string
	}

	// DeadLetterOpts parameters to configure DLQ
	DeadLetterOpts struct {
		QueueName    string
		ExchangeName string
		RoutingKey   string
	}

	// DelayedOpts parameters to configure retry queue exchange
	DelayedOpts struct {
		QueueName    string
		ExchangeName string
		RoutingKey   string
	}

	// Topology used to declare and bind queue, exchanges. Configure dlq and retry
	Topology struct {
		Queue      *QueueOpts
		Exchange   *ExchangeOpts
		Binding    *BindingOpts
		deadLetter *DeadLetterOpts
		delayed    *DelayedOpts
		isBindable bool
	}

	// PUblishOpts
	PublishOpts struct {
		Type      string
		Count     int64
		TraceId   string
		MessageId string
		Delay     time.Duration
	}

	// DeliveryMetadata amqp message received
	DeliveryMetadata struct {
		MessageId string
		XCount    int64
		Type      string
		TraceId   string
		Headers   map[string]interface{}
	}

	// ConsumerHandler
	ConsumerHandler = func(msg any, metadata *DeliveryMetadata) error

	// IRabbitMQMessaging is RabbitMQ  Builder
	IRabbitMQMessaging interface {
		// Declare a new topology
		Declare(opts *Topology) IRabbitMQMessaging

		// Binding bind an exchange/queue with the following parameters without extra RabbitMQ configurations such as Dead Letter.
		ApplyBinds() IRabbitMQMessaging

		// Publish a message
		Publisher(exchange, routingKey string, msg any, opts *PublishOpts) error

		// Create a new goroutine to each dispatcher registered
		//
		// When messages came, some validations will be mad and based on the topology configured message could sent to dql or retry
		Consume() error

		// RegisterDispatcher Add the handler and msg type
		//
		// Each time a message came, we check the queue, and get the available handlers for that queue.
		// After we do a coercion of the msg type to check which handler expect this msg type
		RegisterDispatcher(event string, handler ConsumerHandler, t any) error

		// Build the topology configured
		Build() (IRabbitMQMessaging, error)
	}

	AMQPConnection interface {
		Channel() (*amqp.Channel, error)
	}

	// AMQPChannel is an abstraction for AMQP default channel to improve unit tests
	AMQPChannel interface {
		ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
		ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error
		QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
		QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
		Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
		Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	}

	// Dispatcher struct to register an message handler
	Dispatcher struct {
		Queue         string
		Topology      *Topology
		MsgType       string
		ReflectedType reflect.Value
		Handler       ConsumerHandler
	}

	// IRabbitMQMessaging is the implementation for IRabbitMQMessaging
	RabbitMQMessaging struct {
		Err         error
		logger      logger.ILogger
		conn        AMQPConnection
		ch          AMQPChannel
		config      *env.Configs
		topologies  []*Topology
		dispatchers []*Dispatcher
	}
)

func (d *Topology) ApplyBinds() {
	d.isBindable = true
}

func (d *Topology) RemoveBinds() {
	d.isBindable = false
}
