package rabbitmq

import (
	"reflect"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

type (
	ExchangeKind string

	// QueueOpts declare queue configuration
	QueueOpts struct {
		name           string
		ttl            time.Duration
		retry          *Retry
		withDeadLatter bool
		bindings       []*BindingOpts
	}

	// BindingOpts binds configuration
	BindingOpts struct {
		exchange   string
		routingKey string
	}

	// Retry
	Retry struct {
		NumberOfRetry int64
		DelayBetween  time.Duration
	}

	// ExchangeOpts exchanges to declare
	ExchangeOpts struct {
		name string
		kind ExchangeKind
	}

	Topology interface {
		Exchange(opts *ExchangeOpts) Topology
		FanoutExchanges(exchanges ...string) Topology
		Queue(opts *QueueOpts) Topology
		GetQueueOpts(queue string) *QueueOpts
	}

	// Topology used to declare and bind queue, exchanges. Configure dlq and retry
	topology struct {
		exchanges []*ExchangeOpts
		queues    []*QueueOpts
	}

	// PUblishOpts
	PublishOpts struct {
		Type        string
		Count       int64
		Traceparent string
		MessageId   string
		Delay       time.Duration
	}

	// DeliveryMetadata amqp message received
	DeliveryMetadata struct {
		MessageId   string
		XCount      int64
		Type        string
		Traceparent string
		Headers     map[string]interface{}
	}

	// ConsumerHandler
	ConsumerHandler = func(msg any, metadata *DeliveryMetadata) error

	// IRabbitMQMessaging is RabbitMQ  Builder
	Messaging interface {
		// Declare a new topology
		// Declare(opts *Topology) IRabbitMQMessaging

		// Publish a message
		// Publisher(exchange, routingKey string, msg any, opts *PublishOpts) error

		Channel() AMQPChannel
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

	Dispatcher interface{}

	// Dispatcher struct to register an message handler
	dispatcher struct {
		logger    logging.ILogger
		messaging Messaging
		topology  Topology
		tracer    trace.Tracer

		queues         []string
		msgsTypes      []string
		reflectedTypes []*reflect.Value
		handlers       []ConsumerHandler
	}

	// IRabbitMQMessaging is the implementation for IRabbitMQMessaging
	messaging struct {
		Err      error
		logger   logging.ILogger
		conn     AMQPConnection
		channel  AMQPChannel
		config   *env.Configs
		shotdown chan error
		topology *Topology
		tracer   trace.Tracer
	}
)
