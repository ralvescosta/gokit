# GoKit Tracing Package

The `tracing` package provides distributed tracing capabilities for GoKit-based applications using OpenTelemetry. It enables monitoring and troubleshooting of microservices-based applications by propagating trace context across service boundaries.

## Features

- OpenTelemetry integration for distributed tracing
- OTLP (OpenTelemetry Protocol) exporter support with configurable options
- AMQP (Advanced Message Queuing Protocol) propagation support for tracing across message queues
- Structured logging integration with Zap logger for trace context
- Builder pattern for flexible configuration

## Usage

### Basic Setup with OTLP

```go
import (
    "context"

    "github.com/ralvescosta/gokit/configs"
    "github.com/ralvescosta/gokit/tracing"
)

func main() {
    cfg := configs.NewConfigs() // your application configs

    // Create and configure the OTLP tracing provider
    shutdown, err := tracing.NewOTLP(cfg).
        WithAPIKeyHeader().
        WithTimeout(30 * time.Second).
        Build()

    if err != nil {
        // Handle error
    }

    // Ensure clean shutdown of the tracer provider
    defer shutdown(context.Background())

    // Your application code here with tracing
}
```

### Tracing with AMQP Messages

The package provides utilities for propagating trace context through AMQP messages:

```go
import (
    "github.com/ralvescosta/gokit/tracing"
    "github.com/rabbitmq/amqp091-go"
    "go.opentelemetry.io/otel/trace"
)

// On the producer side
func sendMessage(tracer trace.Tracer, msg amqp.Publishing) {
    ctx, span := tracer.Start(context.Background(), "produce.message")
    defer span.End()

    // Create AMQP header table if it doesn't exist
    if msg.Headers == nil {
        msg.Headers = amqp.Table{}
    }

    // Inject trace context into message headers
    tracing.AMQPPropagator.Inject(ctx, tracing.AMQPHeader(msg.Headers))

    // Send the message with headers containing trace context
}

// On the consumer side
func consumeMessage(tracer trace.Tracer, delivery amqp.Delivery) {
    // Extract trace context and create a new span
    ctx, span := tracing.NewConsumerSpan(tracer, delivery.Headers, "my-queue")
    defer span.End()

    // Process the message with the trace context
    processMessage(ctx, delivery.Body)
}
```

### Adding Trace Information to Logs

The package integrates with the Zap logger to include trace IDs in log entries:

```go
import (
    "github.com/ralvescosta/gokit/logging"
    "github.com/ralvescosta/gokit/tracing"
    "go.uber.org/zap"
)

func processRequest(ctx context.Context, logger logging.Logger) {
    // Add trace information to log entries
    logger.Info("Processing request", tracing.Format(ctx), zap.String("key", "value"))
}
```

## Configuration Options

### ExporterType

The package supports different types of exporters:

- `OTLP_TLS_GRPC_EXPORTER`: Secure OTLP gRPC exporter (with TLS)
- `OTLP_GRPC_EXPORTER`: Insecure OTLP gRPC exporter (without TLS)
- `OTLP_HTTPS_EXPORTER`: OTLP HTTPS exporter
- `JAEGER_EXPORTER`: Jaeger exporter

### OTLP Configuration

The OTLP exporter can be configured with:

- Headers for authentication (including API key)
- Connection timeout
- Reconnection period
- Compression algorithm

## Interfaces

### TracingBuilder

The main interface for configuring and building a tracing provider:

```go
type TracingBuilder interface {
    AddHeader(key, value string) TracingBuilder
    WithHeaders(headers Headers) TracingBuilder
    Type(t ExporterType) TracingBuilder
    Endpoint(s string) TracingBuilder
    Build() (shutdown func(context.Context) error, err error)
}
```

### OTLPTracingBuilder

Extends the TracingBuilder interface with OTLP-specific configuration options:

```go
type OTLPTracingBuilder interface {
    TracingBuilder
    WithAPIKeyHeader() OTLPTracingBuilder
    WithTimeout(t time.Duration) OTLPTracingBuilder
    WithReconnection(t time.Duration) OTLPTracingBuilder
    WithCompression(c OTLPCompression) OTLPTracingBuilder
}
```

## Utility Functions

- `NewConsumerSpan`: Creates a new span for AMQP message consumption
- `Format`: Extracts trace and span IDs from a context for structured logging
- `Message`: Prefixes a string with the package identifier for standardized logging

## Resource Attributes

The tracing provider automatically adds the following resource attributes to all spans:

- `library.language`: "go"
- `service.name`: From application configuration
- `environment`: From application configuration
- `ID`: Current process ID
