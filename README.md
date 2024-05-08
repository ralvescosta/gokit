<h1 align="center">Gokit: Cloud-Native GoLang Toolkit</h1>

<p align="center">
  <a href="https://github.com/ralvescosta/gokit/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  </a>
  <a href="https://goreportcard.com/report/github.com/ralvescosta/gokit">
    <img src="https://goreportcard.com/badge/github.com/ralvescosta/gokit" alt="Go Report Card">
  </a>
</p>

Gokit is a comprehensive GoLang toolkit designed to simplify the development of cloud-native applications. It provides a wide range of features and utilities for logging, observability, messaging, database integration, health checks, migrations, and more. Whether you're a seasoned Go developer or just getting started with cloud-native development, Gokit aims to streamline your workflow and enhance the reliability and efficiency of your applications.

[Table of content]()

  - [Features](#features)
  - [Getting started](#getting-started)
  - [Documentation](#documentation)
  - [License](#license)
  - [Support](#support)

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