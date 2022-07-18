# GoLang gokit 

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_toolkit&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_toolkit)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_toolkit&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_toolkit)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_toolkit&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=ralvescosta_toolkit)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_toolkit&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=ralvescosta_toolkit)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_toolkit&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=ralvescosta_toolkit)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_toolkit&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_toolkit)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_toolkit&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_toolkit)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_toolkit&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_toolkit)

gokit for the must used features in a web application. This project provides severals features, such as: *Environment variables*, *Logging*, *SQL connection management*, *UUID facilities*, *OTL (Open Telemetry)* and *Messaging management*



:warning::construction: **Work In Progress** :construction::warning:

The project is under construction and I want to have a beta version soon.

## gokit 

  - [Environment variables](#environment-variables)
  - [Logging](#logging)
  - [SQL connection management](#sql-connection-management)
  - [UUID facilities](#uuid-facilities)
  - [OTL](#open-telemetry)
  - [Messaging management](#messaging-management)

### Environment variables

- *Package name:* env

### Logging

- *Package name:* logging

### SQL connection management

- *Package name:* sql

### UUID facilities

- *Package name:* uuid

### OTL

- *Package name:* telemetry

- Util Links: [Otel trace](https://opentelemetry.uptrace.dev/instrumentations/?lang=go)

### Messaging Management

- *Package name:* messaging

### Todo

  - [x] Validate the actual implementation and check if the changed Header data will be sent to RabbitMQ
  - [x] Remove Format Log Messages from RabbitMQ pkg and put in Logger pkg
  - [x] Improve errors
  - [x] Change the "logger" pkg name to "logging"
  - [x] Impl Delayed Exchange and Queue
  - [x] Validate Delayed Exchange and Queue
  - [x] Impl the retry strategy
  - [x] Instead of create a new exchange to routing messages to DLQ use the same exchange
  - [x] messaging unit tests
  - [x] uuid unit tests
  - [ ] tracer pkg
  - [ ] create trace-id abstraction for amqp, gRPC and HTTP
  - [ ] adapt messaging to create span in each consumer
  - [ ] adapt sql to create span in each query and based on configuration send the query to the span
