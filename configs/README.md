# Configs

The `configs` package provides a comprehensive configuration framework for GoKit applications. It offers a structured approach to manage different types of configurations like HTTP, database connections, message brokers, authentication, and more.

## Overview

This package contains various configuration structures that can be used to configure different components of a GoKit application. The central `Configs` struct acts as a container for all other configuration types, making it easy to manage and pass around application configuration.

## Key Components

### Core Configuration

- **AppConfigs**: Contains core application settings like environment, application name, logging level, and secret management options.
- **Environment**: A typed enum representing application execution environments (Local, Development, Staging, QA, Production).
- **LogLevel**: Defines logging severity levels (DEBUG, INFO, WARN, ERROR, PANIC).

### Network & API Configuration

- **HTTPConfigs**: Settings for HTTP servers and clients including host, port, and profiling options.
- **IdentityConfigs**: Configuration for identity and authentication services with OAuth2/OIDC support.
- **Auth0Configs**: Auth0-specific authentication settings.

### Database Configuration

- **SQLConfigs**: Database connection settings for SQL databases.
- **DynamoDBConfigs**: Amazon DynamoDB specific configuration.

### Messaging Systems

- **MQTTConfigs**: Configuration for MQTT broker connections with TLS support.
- **RabbitMQConfigs**: Settings for RabbitMQ message broker connections.
- **KafkaConfigs**: Apache Kafka connection and security configuration.

### Cloud Services

- **AWSConfigs**: AWS credential and authentication settings.
- **AWSSecretManagerConfigs**: Configuration for AWS Secret Manager.

### Observability

- **MetricsConfigs**: Settings for application metrics with support for both OpenTelemetry and Prometheus.
- **TracingConfigs**: Distributed tracing configuration using OpenTelemetry.

## Usage Examples

### Basic Application Configuration

```go
appConfig := &configs.AppConfigs{
    GoEnv:   configs.DevelopmentEnv,
    AppName: "my-service",
    LogLevel: configs.INFO,
}
```

### HTTP Server Configuration

```go
httpConfig := &configs.HTTPConfigs{
    Host: "0.0.0.0",
    Port: "8080",
    Addr: "0.0.0.0:8080",
    EnableProfiling: true,
}
```

### Database Configuration

```go
sqlConfig := &configs.SQLConfigs{
    Host:          "localhost",
    Port:          "5432",
    User:          "postgres",
    Password:      "password",
    DbName:        "mydb",
    SecondsToPing: 5,
}
```

### Using the Central Configs Structure

```go
config := &configs.Configs{
    AppConfigs:  appConfig,
    HTTPConfigs: httpConfig,
    SQLConfigs:  sqlConfig,
    // Add other configuration components as needed
}
```

## Environment and Log Level Helpers

The package provides helper functions to convert string representations to typed enums:

```go
// Convert environment string to Environment type
env := configs.NewEnvironment("production") // Returns configs.ProductionEnv

// Convert log level string to LogLevel type
level := configs.NewLogLevel("debug") // Returns configs.DEBUG
```

## Best Practices

1. Use the provided types rather than raw strings for environment and log levels to maintain type safety.
2. Consider using the `configs_builder` package (a companion to this package) for more complex configuration loading from files or environment variables.
3. Store sensitive information like passwords and API keys securely, possibly using the AWS Secret Manager integration provided.
4. Use the central `Configs` struct as a single source of truth for application configuration.

## Related Packages

- `configs_builder`: Helps build configuration objects from various sources like environment variables and files.
- `secrets_manager`: Provides integration with cloud secret management services.
