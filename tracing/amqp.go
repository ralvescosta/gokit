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

type Traceparent struct {
	TraceID    trace.TraceID
	SpanID     trace.SpanID
	TraceFlags trace.TraceFlags
}

var (
	AMQPPropagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
)

type AMQPHeader amqp.Table

func (h AMQPHeader) Set(key, val string) {
	key = strings.ToLower(key)

	h[key] = val
}

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

func (h AMQPHeader) Keys() []string {
	keys := make([]string, 0, len(h))

	for k := range h {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func NewConsumerSpan(tracer trace.Tracer, header amqp.Table, typ string) (context.Context, trace.Span) {
	ctx := AMQPPropagator.Extract(context.Background(), AMQPHeader(header))
	return tracer.Start(ctx, fmt.Sprintf("consume.%s", typ))
}
