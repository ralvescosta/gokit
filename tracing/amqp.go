// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package tracing provides distributed tracing capabilities using OpenTelemetry.
package tracing

import (
	"context"
	"fmt"
	"sort"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Traceparent represents the components of a trace context that are propagated between services
type Traceparent struct {
	TraceID    trace.TraceID
	SpanID     trace.SpanID
	TraceFlags trace.TraceFlags
}

var (
	// AMQPPropagator is a composite propagator that combines TraceContext and Baggage propagation
	// for AMQP messaging contexts
	AMQPPropagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
)

// AMQPHeader wraps amqp.Table to implement the TextMapCarrier interface for OpenTelemetry propagation
type AMQPHeader amqp.Table

// Set sets the value for the given key in the AMQP header (case-insensitive)
func (h AMQPHeader) Set(key, val string) {
	key = strings.ToLower(key)

	h[key] = val
}

// Get retrieves the value for a given key from the AMQP header (case-insensitive)
func (h AMQPHeader) Get(key string) string {
	key = strings.ToLower(key)

	value, ok := h[key]

	if !ok {
		return ""
	}

	toString, ok := value.(string)

	if !ok {
		return ""
	}

	return toString
}

// Keys returns a sorted list of all keys in the AMQP header
func (h AMQPHeader) Keys() []string {
	keys := make([]string, 0, len(h))

	for k := range h {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

// NewConsumerSpan creates a new span for AMQP message consumption with the trace context
// extracted from the message headers
// Parameters:
//   - tracer: The OpenTelemetry tracer to create the span
//   - header: The AMQP message headers containing the trace context
//   - typ: The type of consumer, used to name the span
//
// Returns:
//   - context.Context: Context with the extracted trace information
//   - trace.Span: The new span created for this consumer operation
func NewConsumerSpan(tracer trace.Tracer, header amqp.Table, typ string) (context.Context, trace.Span) {
	ctx := AMQPPropagator.Extract(context.Background(), AMQPHeader(header))
	return tracer.Start(ctx, fmt.Sprintf("consume.%s", typ))
}
