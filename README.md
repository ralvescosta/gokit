# GoLang gokit 

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit&metric=coverage)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit&metric=bugs)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ralvescosta_gokit&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ralvescosta_gokit)

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

- propagation format : 
  - [B3 specification](https://github.com/openzipkin/b3-propagation#single-header): 
    - {trace-id}-{span-id}-{sampling-state}-{parent-span-id}

  - [Jaeger](https://www.jaegertracing.io/docs/1.36/client-libraries/#propagation-format):
    - {trace-id}:{span-id}:0:{flags}
    - {trace-id}:{span-id}:{flags}

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
