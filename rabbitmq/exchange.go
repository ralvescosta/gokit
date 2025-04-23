// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

type (
	// ExchangeKind represents the type of a RabbitMQ exchange.
	// This type defines how messages are routed through the exchange.
	ExchangeKind string

	// ExchangeDefinition represents the configuration of a RabbitMQ exchange.
	// It encapsulates properties like name, durability, auto-delete behavior,
	// exchange type, and additional parameters.
	ExchangeDefinition struct {
		name    string
		durable bool
		delete  bool
		kind    ExchangeKind
		params  map[string]any
	}
)

// String returns the string representation of the ExchangeKind.
func (k ExchangeKind) String() string {
	return string(k)
}

var (
	// FanoutExchange represents a fanout exchange type.
	// Fanout exchanges broadcast all messages to all bound queues.
	FanoutExchange ExchangeKind = "fanout"

	// DirectExchange represents a direct exchange type.
	// Direct exchanges route messages to queues based on a matching routing key.
	DirectExchange ExchangeKind = "direct"
)

// NewDirectExchange creates a new direct exchange definition with the given name.
// Direct exchanges route messages to queues based on exact matching of routing keys.
func NewDirectExchange(name string) *ExchangeDefinition {
	return defaultExchange(name, DirectExchange)
}

// NewFanoutExchange creates a new fanout exchange definition with the given name.
// Fanout exchanges broadcast messages to all queues bound to them.
func NewFanoutExchange(name string) *ExchangeDefinition {
	return defaultExchange(name, FanoutExchange)
}

// NewDirectExchanges creates multiple direct exchange definitions from a list of names.
// This is a convenience function for creating multiple direct exchanges at once.
func NewDirectExchanges(names []string) []*ExchangeDefinition {
	exchanges := []*ExchangeDefinition{}

	for _, name := range names {
		exchanges = append(exchanges, defaultExchange(name, DirectExchange))
	}

	return exchanges
}

// NewFanoutExchanges creates multiple fanout exchange definitions from a list of names.
// This is a convenience function for creating multiple fanout exchanges at once.
func NewFanoutExchanges(names []string) []*ExchangeDefinition {
	exchanges := []*ExchangeDefinition{}

	for _, name := range names {
		exchanges = append(exchanges, defaultExchange(name, FanoutExchange))
	}

	return exchanges
}

// defaultExchange creates a new exchange definition with default settings.
// By default, exchanges are durable and not auto-deleted.
func defaultExchange(name string, kind ExchangeKind) *ExchangeDefinition {
	return &ExchangeDefinition{
		name:    name,
		durable: true,
		delete:  false,
		kind:    kind,
	}
}

// Durable sets the durability flag for the exchange.
// Durable exchanges survive broker restarts.
func (e *ExchangeDefinition) Durable(d bool) *ExchangeDefinition {
	e.durable = d
	return e
}

// Delete sets the auto-delete flag for the exchange.
// Auto-deleted exchanges are removed when no longer in use.
func (e *ExchangeDefinition) Delete(d bool) *ExchangeDefinition {
	e.delete = d
	return e
}

// Params sets additional parameters for the exchange.
// These are passed as arguments when declaring the exchange.
func (e *ExchangeDefinition) Params(p map[string]any) *ExchangeDefinition {
	e.params = p
	return e
}
