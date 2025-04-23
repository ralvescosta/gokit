# Logging

A Go package providing structured logging capabilities built on top of [Zap](https://github.com/uber-go/zap), designed for performance and flexibility across different environments.

## Overview

The logging package provides a simple yet powerful interface for structured logging with different output configurations based on your application's environment. It supports:

- Environment-based configuration (Development, Staging, Production)
- Multiple log levels (Debug, Info, Warn, Error, Fatal)
- Colored console output for development
- JSON formatted logs for production
- File-based logging with concurrent stdout output
- Structured context fields for enhanced log analytics

## Installation

```bash
go get github.com/ralvescosta/gokit/logging
```

## Usage

### Basic Setup

```go
import (
    "github.com/ralvescosta/gokit/configs"
    "github.com/ralvescosta/gokit/logging"
)

func main() {
    // Create app configs
    appConfigs := &configs.Configs{
        AppConfigs: &configs.AppConfigs{
            AppName:  "MyApp",
            GoEnv:    configs.DevelopmentEnv,
            LogLevel: configs.DEBUG,
        },
    }

    // Create a default logger (outputs to stdout)
    logger, err := logging.NewDefaultLogger(appConfigs)
    if err != nil {
        panic(err)
    }

    // Use the logger
    logger.Info("Application started",
        zap.String("version", "1.0.0"),
        zap.Int("port", 8080),
    )
}
```

### File Logging

```go
// Create a file logger (outputs to file and stdout in development)
appConfigs.AppConfigs.LogPath = "./logs/app.log"
logger, err := logging.NewFileLogger(appConfigs)
if err != nil {
    panic(err)
}
```

### Log Levels

The package supports multiple log levels through simple methods:

```go
logger.Debug("Debug message", zap.String("detail", "value"))
logger.Info("Info message", zap.Int("count", 42))
logger.Warn("Warning message", zap.Duration("latency", duration))
logger.Error("Error message", zap.Error(err))
logger.Fatal("Fatal message") // This will exit the application with status code 1
```

### Structured Context

Add structured context to your logs for better filtering and analysis:

```go
logger.Info("User logged in",
    zap.String("user_id", "user123"),
    zap.String("ip_address", "192.168.1.1"),
    zap.String("user_agent", userAgent),
)
```

### Testing

The package includes a MockLogger for use in unit tests:

```go
import (
    "testing"

    "github.com/ralvescosta/gokit/logging"
    "github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
    mockLogger := logging.NewMockLogger()

    // Use the mock logger in your test
    sut := NewSystemUnderTest(mockLogger)

    // Test your system with the mock logger
    // ...
}
```

## Environment-Based Configuration

The logger automatically configures itself based on the environment:

- **Development**: Uses colored console output with a development-friendly format
- **Staging/Production**: Uses JSON formatted logs for better machine parsing

## Performance

By using Zap as the underlying logger, this package inherits Zap's performance benefits:
- Zero allocations during logging
- Minimal CPU overhead
- Type-safe structured logging

## License

MIT
