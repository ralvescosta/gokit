// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

type (
	ExchangeKind string

	ExchangeDefinition struct {
		name    string
		durable bool
		delete  bool
		kind    ExchangeKind
		params  map[string]any
	}
)

func (k ExchangeKind) String() string {
	return string(k)
}

var (
	FanoutExchange ExchangeKind = "fanout"
	DirectExchange ExchangeKind = "direct"
)

func NewDirectExchange(name string) *ExchangeDefinition {
	return defaultExchange(name, DirectExchange)
}

func NewFanoutExchange(name string) *ExchangeDefinition {
	return defaultExchange(name, FanoutExchange)
}

func NewDirectExchanges(names []string) []*ExchangeDefinition {
	exchanges := []*ExchangeDefinition{}

	for _, name := range names {
		exchanges = append(exchanges, defaultExchange(name, DirectExchange))
	}

	return exchanges
}

func NewFanoutExchanges(names []string) []*ExchangeDefinition {
	exchanges := []*ExchangeDefinition{}

	for _, name := range names {
		exchanges = append(exchanges, defaultExchange(name, FanoutExchange))
	}

	return exchanges
}

func defaultExchange(name string, kind ExchangeKind) *ExchangeDefinition {
	return &ExchangeDefinition{
		name:    name,
		durable: true,
		delete:  false,
		kind:    kind,
	}
}

func (e *ExchangeDefinition) Durable(d bool) *ExchangeDefinition {
	e.durable = d
	return e
}

func (e *ExchangeDefinition) Delete(d bool) *ExchangeDefinition {
	e.delete = d
	return e
}

func (e *ExchangeDefinition) Params(p map[string]any) *ExchangeDefinition {
	e.params = p
	return e
}
