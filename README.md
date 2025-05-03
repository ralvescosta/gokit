# Gokit: Cloud-Native Go Toolkit

<p align="center">
  <a href="https://github.com/ralvescosta/gokit/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  </a>
  <a href="https://pkg.go.dev/github.com/ralvescosta/gokit">
    <img src="https://godoc.org/github.com/ralvescosta/gokit?status.svg" alt="Go Doc">
  </a>
  <a href="https://goreportcard.com/report/github.com/ralvescosta/gokit">
    <img src="https://goreportcard.com/badge/github.com/ralvescosta/gokit" alt="Go Report Card">
  </a>
  <a href="https://github.com/ralvescosta/gokit/actions">
    <img src="https://github.com/ralvescosta/gokit/workflows/Go/badge.svg" alt="Build Status">
  </a>
</p>

Gokit is a comprehensive Go toolkit designed to simplify the development of cloud-native applications. It provides a wide range of features and utilities for logging, observability, messaging, database integration, configuration management, and more. Whether you're a seasoned Go developer or just getting started with cloud-native development, Gokit aims to streamline your workflow and enhance the reliability and efficiency of your applications.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Getting Started](#getting-started)
- [Packages](#packages)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Features

- **Configuration Management**: Flexible configuration system with environment-based loading and validation
- **Logging**: Structured logging with multiple output formats and levels
- **Observability**:
  - OpenTelemetry integration for distributed tracing
  - Prometheus metrics for monitoring and analysis
- **Messaging**:
  - MQTT client for IoT and pub/sub messaging
  - RabbitMQ integration for message queuing
  - Kafka support for event streaming
- **Database**:
  - SQL database integration with connection pooling
  - Migration tools for schema management
- **HTTP**: Lightweight HTTP server and client utilities
- **Authentication**: Identity and authorization utilities
- **Secrets Management**: Secure handling of sensitive configuration
- **Utilities**:
  - GUID generation and manipulation
  - Health check endpoints

## Installation

### Prerequisites

Some packages require additional system dependencies:

```bash
# Ubuntu/Debian
sudo apt install libssl-dev build-essential cmake pkg-config llvm-dev libclang-dev clang libmosquitto-dev sqlite3

# CentOS/RHEL
sudo yum install openssl-devel gcc cmake pkg-config llvm-devel clang-devel clang mosquitto-devel sqlite
```

### Go Installation

Add Gokit to your project:

```bash
# Install the entire toolkit
go get -u github.com/ralvescosta/gokit

# Or install specific packages as needed
go get -u github.com/ralvescosta/gokit/logging
go get -u github.com/ralvescosta/gokit/metrics
go get -u github.com/ralvescosta/gokit/mqtt
```

## Getting Started

Here's a simple example of using Gokit's configuration and logging packages:

```go
package main

import (
    "github.com/ralvescosta/gokit/configs_builder"
    "github.com/ralvescosta/gokit/logging"
)

func main() {
    // Build configuration
    cfg, err := configs_builder.NewConfigsBuilder().
        HTTP().
        Metrics().
        Build()
    if err != nil {
        panic(err)
    }

    // Initialize logger
    logger := cfg.Logger
    logger.Info("Application started")

    // Use other Gokit packages...
}
```

## Packages

Gokit is organized into the following packages:

- **auth**: Authentication and authorization utilities
- **configs**: Configuration structures and types
- **configs_builder**: Configuration builder pattern implementation
- **guid**: UUID generation and manipulation
- **httpw**: HTTP wrapper utilities
- **kafka**: Kafka client integration
- **logging**: Structured logging utilities
- **messaging**: Abstractions for messaging brokers like MQTT, RabbitMQ and Kafka
- **metrics**: Prometheus and OpenTelemetry integration.
- **mqtt**: MQTT client implementation
- **rabbitmq**: RabbitMQ client integration
- **secrets_manager**: Secure secrets management
- **sql**: SQL database utilities and migrations
- **tiny_http**: Lightweight HTTP server
- **tracing**: OpenTelemetry tracing integration

## Examples

Check out our [examples repository](https://github.com/ralvescosta/gokit/examples) for complete working examples of Gokit in action, including:

- RabbitMQ publisher and consumer
- HTTP server with metrics and tracing
- Configuration management with environment variables
- Database migrations and queries

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Gokit is released under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

If you have questions, encounter issues, or want to discuss ideas, please:

- Open an issue on the [GitHub Issues](https://github.com/ralvescosta/gokit/issues) page
- Check the [documentation](https://pkg.go.dev/github.com/ralvescosta/gokit) for detailed API information
- Visit our [examples repository](https://github.com/ralvescosta/gokit/examples) for usage examples
