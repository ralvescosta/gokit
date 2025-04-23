# Configs Builder

[![GoDoc](https://godoc.org/github.com/ralvescosta/gokit/configs_builder?status.svg)](https://godoc.org/github.com/ralvescosta/gokit/configs_builder)

The `configs_builder` package provides a fluent interface for building application configurations in Go applications. It simplifies the process of loading configurations from environment variables and `.env` files for various components such as HTTP servers, databases, messaging systems, and more.

## Installation

```bash
go get github.com/ralvescosta/gokit/configs_builder
```

## Features

- Fluent builder pattern for easy configuration setup
- Environment-specific configuration loading (development, staging, production)
- Support for various application components:
  - HTTP server configuration
  - SQL database connections
  - Messaging systems (RabbitMQ, MQTT, Kafka)
  - OpenTelemetry tracing and metrics
  - AWS services (including DynamoDB)
  - Identity/authentication providers

## Usage

### Basic Example

```go
package main

import (
	"github.com/ralvescosta/gokit/configs"
	configsbuilder "github.com/ralvescosta/gokit/configs_builder"
)

func main() {
	// Create a new builder and specify which components to load configuration for
	cfg, err := configsbuilder.NewConfigsBuilder().
		HTTP().       // Load HTTP server configuration
		SQLDatabase(). // Load database configuration
		Tracing().    // Load tracing configuration
		Build()       // Build the final configuration object

	if err != nil {
		panic(err)
	}

	// Use the configured components
	// cfg.HTTPConfigs, cfg.SQLConfigs, etc.
}
```

### Environment Variables

The `configs_builder` reads configuration from environment variables. You can organize these in `.env.development`, `.env.staging`, or `.env.production` files based on your runtime environment.

#### Core Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GO_ENV` | Application environment (development, staging, production) | - |
| `APP_NAME` | Application name | "app" |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | - |
| `LOG_PATH` | Path for log files | "/logs/" |
| `USE_SECRET_MANAGER` | Whether to use secret manager | false |

See the [keys/env_keys.go](keys/env_keys.go) file for a complete list of supported environment variables.

### Configuration Components

#### HTTP Server

```go
cfg, err := configsbuilder.NewConfigsBuilder().
	HTTP().  // Enables HTTP configuration
	Build()
```

Required environment variables:
- `HTTP_HOST` - HTTP server host
- `HTTP_PORT` - HTTP server port

Optional:
- `HTTP_ENABLE_PROFILING` - Enable pprof endpoints (true/false)

#### SQL Database

```go
cfg, err := configsbuilder.NewConfigsBuilder().
	SQLDatabase().  // Enables SQL database configuration
	Build()
```

Required environment variables:
- `SQL_DB_HOST` - Database host
- `SQL_DB_PORT` - Database port
- `SQL_DB_USER` - Database username
- `SQL_DB_PASSWORD` - Database password
- `SQL_DB_NAME` - Database name
- `SQL_DB_SECONDS_TO_PING` - Health check interval in seconds

#### Other Components

The builder supports many more components. For example:

```go
cfg, err := configsbuilder.NewConfigsBuilder().
	RabbitMQ().  // RabbitMQ configuration
	MQTT().      // MQTT configuration
	Metrics().   // Metrics configuration
	Tracing().   // Tracing configuration
	Identity().  // Identity/auth configuration
	Build()
```

## Error Handling

The `configs_builder` validates configuration and returns errors for missing required values:

```go
cfg, err := configsbuilder.NewConfigsBuilder().
	HTTP().
	Build()

if err != nil {
	// Handle missing or invalid configuration
}
```

## Related Packages

- [configs](../configs) - Core configuration structures
- [logging](../logging) - Logging utilities that work with these configurations
- [httpw](../httpw) - HTTP server utilities that use these configurations

## License

MIT License - see [LICENSE](../LICENSE) for details.
