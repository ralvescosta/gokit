// Package kafka provides an implementation of the messaging.Dispatcher interface for Kafka.
// It allows registering handlers for specific message types and sources, and consuming
// messages in a blocking manner, dispatching them to the appropriate handlers.
package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/messaging"
)

// kafkaDispatcher is an implementation of the messaging.Dispatcher interface for Kafka.
// It maintains a registry of handlers for different message types and sources, and ensures
// thread-safe access to the registry using a read-write mutex.
type kafkaDispatcher struct {
	// handlers stores the registered handlers for message types and sources.
	// The outer map key is the source (e.g., Kafka topic), and the inner map key is the message type.
	handlers map[string]map[any]messaging.ConsumerHandler
	mutex    sync.RWMutex
}

// NewDispatcher creates a new instance of kafkaDispatcher.
// It initializes the handlers map and returns a pointer to the dispatcher instance.
func NewDispatcher(configs *configs.Configs) *kafkaDispatcher {
	return &kafkaDispatcher{
		handlers: make(map[string]map[any]messaging.ConsumerHandler),
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

	d.handlers[from][msgType] = handler
	return nil
}

// ConsumeBlocking starts consuming messages and dispatches them to the appropriate registered handlers.
// It iterates over all registered handlers and simulates consuming messages from Kafka.
// For each message, it invokes the corresponding handler with the message and its metadata.
func (d *kafkaDispatcher) ConsumeBlocking() {
	for {
		d.mutex.RLock()
		for from, handlersForSource := range d.handlers {
			for msgType, handler := range handlersForSource {
				// Simulate consuming a message from Kafka.
				msg := fmt.Sprintf("message for type %v from source %s", msgType, from)
				metadata := fmt.Sprintf("metadata for type %v from source %s", msgType, from)

				ctx := context.Background()
				if err := handler(ctx, msg, metadata); err != nil {
					fmt.Printf("Error handling message: %v\n", err)
				}
			}
		}
		d.mutex.RUnlock()
	}
}
