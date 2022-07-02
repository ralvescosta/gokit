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

	Retry struct {
		NumberOfRetry int64
		DelayBetween  time.Duration
	}

	QueueOpts struct {
		Name           string
		TTL            time.Duration
		Retryable      *Retry
		WithDeadLatter bool
	}

	ExchangeOpts struct {
		Name     string
		Type     ExchangeKind
		Bindings []string
	}

	BindingOpts struct {
		RoutingKey        string
		dlqRoutingKey     string
		delayedRoutingKey string
	}

	DeadLetterOpts struct {
		QueueName    string
		ExchangeName string
		RoutingKey   string
	}

	DelayedOpts struct {
		QueueName    string
		ExchangeName string
		RoutingKey   string
	}

	Topology struct {
		Queue      *QueueOpts
		Exchange   *ExchangeOpts
		Binding    *BindingOpts
		deadLetter *DeadLetterOpts
		delayed    *DelayedOpts
		isBindable bool
	}

	PublishOpts struct {
		Type      string
		Count     int64
		TraceId   string
		MessageId string
		Delay     time.Duration
	}

	DeliveryMetadata struct {
		MessageId string
		XCount    int64
		Type      string
		TraceId   string
		Headers   map[string]interface{}
	}

	ConsumerHandler = func(msg any, metadata *DeliveryMetadata) error

	// IRabbitMQMessaging is RabbitMQ  Builder
	IRabbitMQMessaging interface {
		// Declare
		Declare(opts *Topology) IRabbitMQMessaging

		// Binding bind an exchange/queue with the following parameters without extra RabbitMQ configurations such as Dead Letter.
		ApplyBinds() IRabbitMQMessaging

		Publisher(exchange, routingKey string, msg any, opts *PublishOpts) error

		Consume() error

		// AddDispatcher Add the handler and msg type
		//
		// Each time a message came, we check the queue, and get the available handlers for that queue.
		// After we do a coercion of the msg type to check which handler expect this msg type
		RegisterDispatcher(event string, handler ConsumerHandler, t any) error

		Build() (IRabbitMQMessaging, error)
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
		config      *env.Configs
		conn        *amqp.Connection
		ch          AMQPChannel
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
