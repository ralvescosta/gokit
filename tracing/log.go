// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type traceLog struct {
	TraceID string
	SpanID  string
}

func (u *traceLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("trace_id", u.TraceID)
	enc.AddString("span_id", u.SpanID)
	return nil
}

func Format(ctx context.Context) zapcore.Field {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return zap.Skip()
	}

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	return zap.Inline(&traceLog{traceID, spanID})
}
