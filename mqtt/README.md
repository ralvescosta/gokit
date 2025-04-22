# MQTT Package Documentation

The `mqtt` package provides a robust and idiomatic implementation for working with MQTT brokers in Go. It includes features for publishing, subscribing, and managing MQTT connections, while adhering to Go best practices.

## Features

- **Client Management**: Establish and manage connections to MQTT brokers.
- **Publishing**: Publish messages to topics with or without context deadlines.
- **Subscription Management**: Register and consume messages from topics.
- **Error Handling**: Comprehensive error handling for common MQTT operations.
- **Tracing**: Integrated with OpenTelemetry for distributed tracing.

## Installation

To use the `mqtt` package, add it to your project:

```bash
go get github.com/ralvescosta/gokit/mqtt
```

## Usage

### Creating an MQTT Client

```go
import (
	"github.com/ralvescosta/gokit/mqtt"
    configsBuilder "github.com/ralvescosta/gokit/configs_builder"
)

func main() {
	cfgs, err := configsBuilder.
		NewConfigsBuilder().
		MQTT().
		Build()
	if err != nil {
		cfgs.Logger.Fatal(err.Error())
	}

	client := mqtt.NewMQTTClient(cfgs)

	if err := client.Connect(); err != nil {
		panic(err)
	}

	defer client.Client().Disconnect(250)
}
```

### Publishing Messages

#### With Context Deadline

```go
import (
	"context"
	"github.com/ralvescosta/gokit/mqtt"
)

func publishWithContext(publisher mqtt.Publisher) {
	ctx := context.Background()
	err := publisher.PubCtx(ctx, "example/topic", mqtt.AtLeastOnce, "Hello, MQTT!")
	if err != nil {
		panic(err)
	}
}
```

#### Without Context Deadline

```go
func publishWithoutContext(publisher mqtt.Publisher) {
	err := publisher.Pub("example/topic", mqtt.AtLeastOnce, "Hello, MQTT!")
	if err != nil {
		panic(err)
	}
}
```

### Subscribing to Topics

```go
import (
	"os"
	"github.com/ralvescosta/gokit/mqtt"
)

func subscribe(dispatcher mqtt.Dispatcher) {
	handler := func(ctx context.Context, topic string, qos mqtt.QoS, payload []byte) error {
		fmt.Printf("Received message on topic %s: %s\n", topic, string(payload))
		return nil
	}

	err := dispatcher.Register("example/topic", mqtt.AtLeastOnce, handler)
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	dispatcher.ConsumeBlocking(stop)
}
```

## Error Handling

The `mqtt` package provides predefined errors for common issues:

- `ConnectionFailureError`: Indicates a failure to connect to the MQTT broker.
- `EmptyTopicError`: Indicates that the topic for a subscription cannot be empty.
- `NillHandlerError`: Indicates that the handler for a subscription cannot be nil.
- `NillPayloadError`: Indicates that the payload for a publish operation cannot be nil.
- `InvalidQoSError`: Indicates that the provided QoS value is invalid.

## Tracing

The `mqtt` package integrates with OpenTelemetry for distributed tracing. Each message handler creates a new span with the topic name as the span name.

## License

This package is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details.

## Contributing

Contributions are welcome! Please read the [contributing guidelines](../CONTRIBUTING.md) before submitting a pull request.

## Support

For support, please open an issue in the repository or contact the maintainers.