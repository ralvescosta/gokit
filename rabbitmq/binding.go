// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

type (
	// ExchangeBindingDefinition represents a binding between two exchanges.
	// It defines how messages are routed from a source exchange to a destination exchange
	// based on a routing key and optional arguments.
	ExchangeBindingDefinition struct {
		source      string
		destination string
		routingKey  string
		args        map[string]interface{}
	}

	// QueueBindingDefinition represents a binding between an exchange and a queue.
	// It defines how messages are routed from an exchange to a queue
	// based on a routing key and optional arguments.
	QueueBindingDefinition struct {
		routingKey string
		queue      string
		exchange   string
		args       map[string]interface{}
	}
)

// NewExchangeBiding creates a new exchange binding definition.
// This defines how messages are routed between exchanges.
func NewExchangeBiding() *ExchangeBindingDefinition {
	return &ExchangeBindingDefinition{}
}

// NewQueueBinding creates a new queue binding definition.
// This defines how messages are routed from an exchange to a queue.
func NewQueueBinding() *QueueBindingDefinition {
	return &QueueBindingDefinition{}
}

// RoutingKey sets the routing key for this queue binding.
// The routing key is used to filter messages from the exchange to the queue.
func (b *QueueBindingDefinition) RoutingKey(key string) *QueueBindingDefinition {
	b.routingKey = key
	return b
}

// Queue sets the queue name for this binding.
// This is the destination queue that will receive messages from the exchange.
func (b *QueueBindingDefinition) Queue(name string) *QueueBindingDefinition {
	b.queue = name
	return b
}

// Exchange sets the exchange name for this binding.
// This is the source exchange from which messages will be routed to the queue.
func (b *QueueBindingDefinition) Exchange(name string) *QueueBindingDefinition {
	b.exchange = name
	return b
}
