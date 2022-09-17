package tracing

import (
	"context"
	"strings"

	"github.com/ralvescosta/gokit/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Traceparent struct {
	TraceID    trace.TraceID
	SpanID     trace.SpanID
	TraceFlags trace.TraceFlags
}

func TraceparentFromString(traceparent string) (*Traceparent, error) {
	v := strings.Split(traceparent, "-")

	if len(v) < 3 {
		return nil, errors.ErrorAMQPBadTraceparent
	}

	traceID := [16]byte{}
	copy(traceID[:], v[0])

	spanID := [8]byte{}
	copy(spanID[:], v[1])

	return &Traceparent{
		TraceID:    trace.TraceID(traceID),
		SpanID:     trace.SpanID(spanID),
		TraceFlags: trace.TraceFlags([]byte(v[2])[0]),
	}, nil
}

func ContextFromTraceparent(ctx context.Context, traceparent string) (context.Context, error) {
	parent, err := TraceparentFromString(traceparent)
	if err != nil {
		return nil, err
	}

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    parent.TraceID,
		SpanID:     parent.SpanID,
		TraceFlags: parent.TraceFlags,
		Remote:     true,
	})

	return trace.ContextWithRemoteSpanContext(ctx, spanCtx), nil
}

func SpanFromAMQPTraceparent(tracer trace.Tracer, traceparent, name, exch, queue string) (context.Context, trace.Span, error) {
	ctx := context.Background()

	ctx, err := ContextFromTraceparent(ctx, traceparent)
	if err != nil {
		return ctx, nil, err
	}

	ctx, span := tracer.Start(ctx, name)

	span.SetAttributes(attribute.KeyValue{
		Key:   "amqp.exchange",
		Value: attribute.StringValue(exch),
	})
	span.SetAttributes(attribute.KeyValue{
		Key:   "amqp.queue",
		Value: attribute.StringValue(queue),
	})

	return ctx, span, nil
}

func StringTraceparentFromCtx(ctx context.Context) {
	println(ctx)
}
