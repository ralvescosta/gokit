// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package kafka

import (
	"context"
	"errors"
	"sync"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/messaging"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// kafkaDispatcher is an implementation of the messaging.Dispatcher interface for Kafka.
// It maintains a registry of handlers for different message types and sources, and ensures
// thread-safe access to the registry using a read-write mutex.
type kafkaDispatcher struct {
	logger logging.Logger

	// handlers stores the registered handlers for message types and sources.
	// The outer map key is the source (e.g., Kafka topic), and the inner map key is the message type.
	handlers map[string]map[any]messaging.ConsumerHandler
	mutex    sync.RWMutex

	// kafkaReaders is a slice of Kafka readers used to consume messages.
	kafkaReaders []*kafka.Reader
}

// NewDispatcher creates a new instance of kafkaDispatcher.
// It initializes the handlers map and returns a pointer to the dispatcher instance.
func NewDispatcher(configs *configs.Configs) *kafkaDispatcher {
	return &kafkaDispatcher{
		logger:       configs.Logger,
		handlers:     make(map[string]map[any]messaging.ConsumerHandler),
		kafkaReaders: []*kafka.Reader{},
	}
}

// Register associates a message type and source with a specific messaging.ConsumerHandler.
// It ensures that the same handler is not registered multiple times for the same message type and source.
//
// Parameters:
// - from: The source of the message (e.g., Kafka topic).
// - msgType: The type of the message.
// - handler: The handler function to process the message.
//
// Returns:
// - An error if a handler is already registered for the given message type and source.
func (d *kafkaDispatcher) Register(from string, msgType any, handler messaging.ConsumerHandler) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, exists := d.handlers[from]; !exists {
		d.handlers[from] = make(map[any]messaging.ConsumerHandler)
	}

	if _, exists := d.handlers[from][msgType]; exists {
		return errors.New("handler already registered for this message type and source")
	}

	// Create a new Kafka reader for the given source (topic) and add it to the slice.
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"}, // Replace with actual broker addresses
		GroupID: "default-group",            // Replace with actual group ID
		Topic:   from,
	})

	d.kafkaReaders = append(d.kafkaReaders, reader)
	d.handlers[from][msgType] = handler

	return nil
}

// ConsumeBlocking starts consuming messages from Kafka and dispatches them to the appropriate registered handlers.
// It creates a separate goroutine for each Kafka reader to consume messages concurrently.
func (d *kafkaDispatcher) ConsumeBlocking() {
	for _, reader := range d.kafkaReaders {
		go func(r *kafka.Reader) {
			for {
				msg, err := r.ReadMessage(context.Background())
				if err != nil {
					d.logger.Error("Error reading message from Kafka", zap.Error(err))
					continue
				}

				d.mutex.RLock()
				handlersForSource, exists := d.handlers[msg.Topic]
				d.mutex.RUnlock()

				if !exists {
					d.logger.Warn("No handlers registered for topic", zap.String("topic", msg.Topic))
					continue
				}

				handler, exists := handlersForSource[string(msg.Key)]
				if !exists {
					d.logger.Warn("No handler registered for message type", zap.String("messageType", string(msg.Key)))
					continue
				}

				ctx := context.Background()
				if err := handler(ctx, msg.Value, msg.Headers); err != nil {
					d.logger.Error("Error handling message", zap.Error(err))
				}
			}
		}(reader)
	}
}
