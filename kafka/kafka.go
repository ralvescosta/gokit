// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package kafka implements a robust interface for working with Apache Kafka messaging systems.
//
// This package provides implementations of the messaging interfaces for the Kafka message broker,
// allowing applications to publish messages to Kafka topics and consume messages with typed
// handlers. It wraps the segmentio/kafka-go library to provide a more ergonomic API that
// integrates with the broader GoKit ecosystem.
//
// # Features
//
// - Type-safe message publishing with the Publisher interface
// - Consumer-side message dispatching with the Dispatcher interface
// - Support for message headers and metadata
// - Integration with GoKit's configuration, logging, and tracing systems
// - Error handling and recovery mechanisms
//
// # Components
//
// - Publisher: Sends messages to Kafka topics with support for deadlines and options
// - Dispatcher: Consumes messages from topics and routes them to registered handlers
//
// # Usage
//
// Publisher example:
//
//	publisher := kafka.NewPublisher(configs)
//	topic := "orders"
//	key := "new-order"
//	err := publisher.Publish(ctx, &topic, nil, &key, orderData)
//
// Dispatcher example:
//
//	dispatcher := kafka.NewDispatcher(configs)
//	err := dispatcher.Register("orders", OrderCreated{}, handleOrderCreated)
//	dispatcher.ConsumeBlocking()
//
// # Configuration
//
// The package uses the GoKit configs package for configuration:
//
//	// Required Kafka configuration
//	configs.KafkaConfigs.Host = "localhost:9092"
//
// See the configs package documentation for all available configuration options.
package kafka
