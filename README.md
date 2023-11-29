<h1 align="center">Gokit: Cloud-Native GoLang Toolkit</h1>

<p align="center">
  <a href="https://github.com/ralvescosta/gokit/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  </a>
  <a href="https://goreportcard.com/report/github.com/ralvescosta/gokit">
    <img src="https://goreportcard.com/badge/github.com/ralvescosta/gokit" alt="Go Report Card">
  </a>
</p>

:warning::construction: **Work In Progress** :construction::warning:

Gokit is a comprehensive GoLang toolkit designed to simplify the development of cloud-native applications. It provides a wide range of features and utilities for logging, observability, messaging, database integration, health checks, migrations, and more. Whether you're a seasoned Go developer or just getting started with cloud-native development, Gokit aims to streamline your workflow and enhance the reliability and efficiency of your applications.

[Table of content]()

  - [Features](#features)
  - [Getting started](#getting-started)
  - [Documentation](#documentation)
  - [License](#license)
  - [Support](#support)
  - [Sonar Metrics](#sonar-metrics)

# Features

- **Logging:** Streamline structured logging in your Go applications for improved debugging and monitoring capabilities.

- **OpenTelemetry and Tracing:** Seamlessly integrate with OpenTelemetry to enable distributed tracing and gain insights into your application's performance.

- **OpenTelemetry and Prometheus Metrics:** Integrate with Prometheus to collect and expose application metrics for monitoring and analysis.

- **Messaging:** Support for MQTT and RabbitMQ, simplifying the implementation of messaging patterns in your applications.

- **Database:** Integration with PostgreSQL for efficient data storage and retrieval.

- **Health Checks:** Implement health checks to ensure the reliability of your services and applications.

- **Migrator:** Simplify database schema migrations to manage changes in your data models.

## Getting Started

1- **Installing Linux packages:**  Some crates require some additional packages to work pronely

```bash
sudo apt install libssl-dev build-essential cmake pkg-config llvm-dev libclang-dev clang libmosquitto-dev sqlite3
```

2- **Add in our project** To use one of these crates add one of the packages in your project:

```bash
go get -u github.com/ralvescosta/gokit/logging
```

# Documentation

For detailed documentation and usage examples, please visit our [gokit example repository](https://github.com/ralvescosta/gokit_examples)

# License

Ruskit is released under the MIT License. See [LICENSE](https://github.com/ralvescosta/gokit/blob/main/LICENSE) for more details.

# Support

If you have questions, encounter issues, or want to discuss ideas, please open an issue on the [GitHub Issues](https://github.com/ralvescosta/gokit/issues) page.

# Sonar Metrics

|   |   |   |   |   |
|---|---|---|---|---|
| [configs](https://github.com/ralvescosta/gokit/tree/main/configs) | [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_configs&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_configs) | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_configs&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_configs) | [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_configs&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_configs) | [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_configs&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_configs) | [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_configs&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_configs) | [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_configs&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_configs) |
| [guid](https://github.com/ralvescosta/gokit/tree/main/guid) | [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_guid&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_guid) | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_guid&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_guid) | [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_guid&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_guid) | [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_guid&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_guid) | [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_guid&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_guid) | [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_guid&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_env) |
| [httpw](https://github.com/ralvescosta/gokit/tree/main/httpw) | [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_httpw&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_httpw) | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_httpw&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_httpw) | [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_httpw&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_httpw) | [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_httpw&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_httpw) | [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_httpw&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_httpw) | [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_httpw&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_httpw) |
| [logging](https://github.com/ralvescosta/gokit/tree/main/logging) | [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_logging&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_logging) | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_logging&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_logging) | [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_logging&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_logging) | [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_logging&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_logging) | [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_logging&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_logging) | [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_logging&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_logging) |
| [metrics](https://github.com/ralvescosta/gokit/tree/main/metrics) | [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_metrics&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_metrics) | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_metrics&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_metrics) | [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_metrics&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_metrics) | [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_metrics&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_metrics) | [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_metrics&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_metrics) | [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_metrics&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_metrics) |
| [rabbitmq](https://github.com/ralvescosta/gokit/tree/main/rabbitmq) | [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_rabbitmq&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_rabbitmq) | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_rabbitmq&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_rabbitmq) | [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_rabbitmq&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_rabbitmq) | [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_rabbitmq&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_rabbitmq) | [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_rabbitmq&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_rabbitmq) | [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_rabbitmq&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_rabbitmq) |
| [secrets_manager](https://github.com/ralvescosta/gokit/tree/main/secrets_manager) | [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_secrets_manager&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_secrets_manager) | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_secrets_manager&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_secrets_manager) | [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_secrets_manager&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_secrets_manager) | [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_secrets_manager&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_secrets_manager) | [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_secrets_manager&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_secrets_manager) | [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_secrets_manager&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_secrets_manager) |
| [sql](https://github.com/ralvescosta/gokit/tree/main/sql) | [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_sql&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_sql) | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_sql&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_sql) | [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_sql&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_sql) | [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_sql&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_sql) | [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_sql&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_sql) | [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_sql&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_sql) |
| [tracing](https://github.com/ralvescosta/gokit/tree/main/tracing) | [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_tracing&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_tracing) | [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_tracing&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_tracing) | [![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_tracing&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_tracing) | [![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_tracing&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_tracing) | [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_tracing&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_tracing) | [![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit_tracing&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit_tracing) |
| [mqtt]() | [Quality Gate Status]() | [Lines of Code]() | [Vulnerabilities]() | [Bugs]() | [Security Rating]() |




## Todo
- [] Sql package
  - [] remove graceful shotdown
  - [] Graceful shotdown strategy

- [] Improve gokit examples