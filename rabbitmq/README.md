# RabbitMQ Package

The `rabbitmq` package provides a robust and idiomatic implementation for working with RabbitMQ in Go applications. It includes features for connection management, topology setup, publishing, and consuming messages with support for advanced patterns like dead-letter queues and retry mechanisms.

## Features

- **Connection Management**: Establish and manage connections to RabbitMQ brokers
- **Topology Definition**: Define exchanges, queues, and bindings using a fluent API
- **Publishing**: Publish messages to exchanges or directly to queues
- **Message Consumption**: Register handlers for message processing with automatic deserialization
- **Error Handling**: Comprehensive error handling with custom error types
- **Dead Letter Queues**: Support for DLQ pattern for failed message handling
- **Retry Mechanism**: Configurable retry mechanism for transient failures
- **Tracing**: Integration with OpenTelemetry for distributed tracing

## Installation

To use the `rabbitmq` package, add it to your project:

```bash
go get github.com/ralvescosta/gokit/rabbitmq
```

## Usage

### Establishing a Connection

```go
import (
	"github.com/ralvescosta/gokit/configs_builder"
	"github.com/ralvescosta/gokit/rabbitmq"
)

func main() {
	// Create configurations
	cfgs, err := configsBuilder.
		NewConfigsBuilder().
		RabbitMQ().
		Build()
	if err != nil {
		panic(err)
	}

	// Establish connection and create channel
	conn, ch, err := rabbitmq.NewConnection(cfgs)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
```

### Defining Topology

```go
// Create topology manager
topology := rabbitmq.NewTopology(cfgs).Channel(ch)

// Define exchanges
fanoutExchange := rabbitmq.NewFanoutExchange("notifications")
directExchange := rabbitmq.NewDirectExchange("orders")
topology.Exchange(fanoutExchange).Exchange(directExchange)

// Define queues
ordersQueue := rabbitmq.NewQueue("orders").WithDQL().WithRetry(time.Second*5, 3)
notificationsQueue := rabbitmq.NewQueue("notifications")
topology.Queue(ordersQueue).Queue(notificationsQueue)

// Define bindings
orderBinding := rabbitmq.NewQueueBinding().
	Queue("orders").
	Exchange("orders").
	RoutingKey("new_order")
topology.QueueBinding(orderBinding)

// Apply topology to RabbitMQ
if err := topology.Apply(); err != nil {
	panic(err)
}
```

### Publishing Messages

```go
// Create publisher
publisher := rabbitmq.NewPublisher(cfgs, ch)

// Create message
type OrderCreated struct {
	OrderID string `json:"order_id"`
	Amount  float64 `json:"amount"`
}
msg := OrderCreated{OrderID: "123", Amount: 99.99}

// Publish to exchange with routing key
ctx := context.Background()
if err := publisher.Publish(ctx, "orders", "new_order", msg); err != nil {
	log.Printf("Failed to publish: %v", err)
}

// Publish directly to queue
if err := publisher.SimplePublish(ctx, "orders", msg); err != nil {
	log.Printf("Failed to publish: %v", err)
}
```

### Consuming Messages

```go
// Create dispatcher
dispatcher := rabbitmq.NewDispatcher(cfgs, ch, topology.GetQueuesDefinition())

// Define message type
type OrderCreated struct {
	OrderID string `json:"order_id"`
	Amount  float64 `json:"amount"`
}

// Register message handler
err := dispatcher.Register("orders", OrderCreated{}, func(ctx context.Context, msg any, metadata any) error {
	order, ok := msg.(*OrderCreated)
	if !ok {
		return fmt.Errorf("invalid message type")
	}

	fmt.Printf("Processing order: %s, amount: %.2f\n", order.OrderID, order.Amount)
	return nil
})
if err != nil {
	panic(err)
}

// Start consuming messages
dispatcher.ConsumeBlocking()
```

## Error Handling

The package provides predefined errors for common issues:

- `NullableChannelError`: Returned when a channel operation is attempted on a nil channel
- `NotFoundQueueDefinitionError`: Returned when a queue definition cannot be found
- `InvalidDispatchParamsError`: Returned when invalid parameters are provided to a dispatch operation
- `QueueDefinitionNotFoundError`: Returned when no queue definition is found for a specified queue
- `ReceivedMessageWithUnformattedHeaderError`: Returned when a message has incorrectly formatted headers
- `RetryableError`: Indicates that a message processing failed but can be retried later

## Advanced Features

### Dead Letter Queues

```go
// Create a queue with Dead Letter Queue support
ordersQueue := rabbitmq.NewQueue("orders").WithDQL()

// Any rejected or failed messages will be routed to "orders-dlq"
```

### Retry Mechanism

```go
// Create a queue with retry mechanism
// Messages will be retried up to 3 times with 5-second delay between attempts
ordersQueue := rabbitmq.NewQueue("orders").WithRetry(time.Second*5, 3)
```

### Tracing

The `rabbitmq` package automatically integrates with OpenTelemetry to provide distributed tracing for message publishing and consumption. Trace context is propagated through message headers.

## License

This package is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details.
