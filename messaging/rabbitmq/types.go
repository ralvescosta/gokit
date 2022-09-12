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

	TopologyBuilder interface {
		Exchange(opts *ExchangeOpts) TopologyBuilder
		FanoutExchanges(exchanges ...string) TopologyBuilder
		Queue(opts *QueueOpts) TopologyBuilder
	}

	// Topology used to declare and bind queue, exchanges. Configure dlq and retry
	Topology struct {
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
	IRabbitMQMessaging interface {
		// Declare a new topology
		// Declare(opts *Topology) IRabbitMQMessaging

		// Binding bind an exchange/queue with the following parameters without extra RabbitMQ configurations such as Dead Letter.
		// ApplyBinds() IRabbitMQMessaging

		// Publish a message
		// Publisher(exchange, routingKey string, msg any, opts *PublishOpts) error

		// Create a new goroutine to each dispatcher registered
		//
		// When messages came, some validations will be mad and based on the topology configured message could sent to dql or retry
		// Consume() error

		// RegisterDispatcher Add the handler and msg type
		//
		// Each time a message came, we check the queue, and get the available handlers for that queue.
		// After we do a coercion of the msg type to check which handler expect this msg type
		// RegisterDispatcher(event string, handler ConsumerHandler, t any) error

		// Build the topology configured
		// Build() (IRabbitMQMessaging, error)
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
		queues         []string
		msgsTypes      []string
		reflectedTypes []*reflect.Value
		handlers       []ConsumerHandler
		messaging      *RabbitMQMessaging
	}

	// IRabbitMQMessaging is the implementation for IRabbitMQMessaging
	RabbitMQMessaging struct {
		Err         error
		logger      logging.ILogger
		conn        AMQPConnection
		channel     AMQPChannel
		config      *env.Configs
		shotdown    chan error
		topology    *Topology
		dispatchers []*Dispatcher
		tracer      trace.Tracer
	}
)
