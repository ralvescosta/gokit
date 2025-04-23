# GoKit - Metrics

The metrics package provides tools for collecting, exporting, and monitoring metrics with support for multiple exporters like OpenTelemetry Protocol (OTLP) and Prometheus.

## Overview

The `metrics` package is part of the GoKit framework and offers a comprehensive solution for application monitoring through metrics collection. It supports:

- OpenTelemetry Protocol (OTLP) metrics exporter with gRPC transport
- Prometheus metrics exporter
- HTTP middleware for collecting request metrics
- System metrics collectors for Go runtime statistics (memory, goroutines, GC, etc.)

## Installation

```bash
go get github.com/ralvescosta/gokit/metrics
```

## Usage

### OpenTelemetry Protocol (OTLP) Exporter

The OTLP exporter allows you to send metrics to any OpenTelemetry-compatible backend like Jaeger or a custom collector.

```go
import (
    "context"
    "time"
    
    "github.com/ralvescosta/gokit/configs"
    "github.com/ralvescosta/gokit/metrics"
)

func setupOTLPMetrics(cfgs *configs.Configs) (shutdown func(context.Context) error, err error) {
    // Create a new OTLP builder
    builder := metrics.NewOTLPBuilder(cfgs)
    
    // Configure OTLP exporter
    return builder.
        WithAPIKeyHeader().                       // Add API key header if configured
        Endpoint("localhost:4317").               // Set OTLP endpoint
        WithTimeout(30 * time.Second).            // Set timeout
        WithCompression(metrics.OtlpGzipCompressions). // Set compression
        Provider()                                // Create the provider
}
```

### Prometheus Exporter

The Prometheus exporter allows your application to expose metrics for scraping by Prometheus.

```go
import (
    "context"
    "net/http"
    
    "github.com/ralvescosta/gokit/configs"
    "github.com/ralvescosta/gokit/metrics"
)

func setupPrometheusMetrics(cfgs *configs.Configs) (http.Handler, func(context.Context) error, error) {
    // Create a new Prometheus exporter
    prom := metrics.NewPrometheus(cfgs)
    
    // Get the HTTP handler for metrics endpoint
    handler := prom.HTTPHandler()
    
    // Configure the provider
    shutdown, err := prom.Provider()
    if err != nil {
        return nil, nil, err
    }
    
    return handler, shutdown, nil
}

// Use the handler in your HTTP server
func exposeMetricsEndpoint(handler http.Handler) {
    http.Handle("/metrics", handler)
    http.ListenAndServe(":8080", nil)
}
```

### HTTP Metrics Middleware

The HTTP middleware collects metrics about request counts and durations.

```go
import (
    "net/http"
    
    httpMetrics "github.com/ralvescosta/gokit/metrics/http"
)

func setupHTTPServer() {
    // Create the metrics middleware
    middleware, err := httpMetrics.NewHTTPMetricsMiddleware()
    if err != nil {
        panic(err)
    }
    
    // Create your HTTP handler
    myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    // Wrap your handler with the metrics middleware
    http.Handle("/", middleware.Handler(myHandler))
    http.ListenAndServe(":8080", nil)
}
```

### System Metrics Collection

The system package provides collectors for Go runtime metrics like memory usage, garbage collection stats, and active goroutines.

```go
import (
    "github.com/ralvescosta/gokit/logging"
    "github.com/ralvescosta/gokit/metrics/system"
)

func setupSystemMetrics(logger logging.Logger) error {
    // Initialize basic metrics collectors (memory and system)
    return system.BasicMetricsCollector(logger)
}
```

## Metrics Types

The package includes several metrics types based on OpenTelemetry:

1. **Counters**: For values that only go up (like request counts)
2. **Gauges**: For values that can go up and down (like memory usage)
3. **Histograms**: For measurements like request durations with statistical analysis

## Best Practices

1. **Initialization**: Initialize metrics early in your application lifecycle
2. **Cleanup**: Use the shutdown functions returned by providers to properly cleanup resources
3. **Context**: Pass context to metrics operations for proper propagation
4. **Naming**: Use consistent naming conventions for your metrics
5. **Labels/Attributes**: Keep cardinality low by limiting the number of unique values

## Integration with Other GoKit Packages

The metrics package integrates with other GoKit packages:

- Works with the `configs` package for configuration
- Uses the `logging` package for logging metrics events
- Can be used alongside the `tracing` package for complete observability

## License

MIT License - See the LICENSE file for details.