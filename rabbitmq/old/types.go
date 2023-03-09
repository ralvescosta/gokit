package rabbitmq

import (
	"context"
	"os"
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
		dqlName        string
		retryName      string
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
	topologyImpl struct {
		exchanges []*ExchangeOpts
		queues    []*QueueOpts
	}

	// PUblishOpts
	PublishOpts struct {
		Type      string
		Count     int64
		MessageId string
		Delay     time.Duration
	}

	// DeliveryMetadata amqp message received
	DeliveryMetadata struct {
		MessageId string
		XCount    int64
		Type      string
		Headers   map[string]interface{}
	}

	// ConsumerHandler
	ConsumerHandler = func(ctx context.Context, msg any, metadata *DeliveryMetadata) error

	// IRabbitMQMessaging is RabbitMQ  Builder
	Messaging interface {
		Channel() AMQPChannel
		InstallTopology(topology Topology) (Messaging, error)
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

	Dispatcher interface {
		RegisterDispatcher(queue string, msg any, handler ConsumerHandler) error
		ConsumeBlocking(ch chan os.Signal)
	}

	// Dispatcher struct to register an message handler
	dispatcherImpl struct {
		logger    logging.Logger
		messaging Messaging
		topology  Topology
		tracer    trace.Tracer

		queues         []string
		msgsTypes      []string
		reflectedTypes []*reflect.Value
		handlers       []ConsumerHandler
	}

	// IRabbitMQMessaging is the implementation for IRabbitMQMessaging
	messagingImpl struct {
		Err     error
		logger  logging.Logger
		conn    AMQPConnection
		channel AMQPChannel
		config  *env.RabbitMQConfigs
		tracer  trace.Tracer
	}
)
