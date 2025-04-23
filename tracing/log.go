// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package tracing provides distributed tracing capabilities using OpenTelemetry.
package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// traceLog is a struct that holds trace and span IDs for structured logging
type traceLog struct {
	TraceID string
	SpanID  string
}

// MarshalLogObject implements zapcore.ObjectMarshaler interface for traceLog
// This allows the trace information to be added to zap logs in a structured way
func (u *traceLog) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("trace_id", u.TraceID)
	enc.AddString("span_id", u.SpanID)
	return nil
}

// Format extracts trace and span IDs from a context and returns them as a zap field
// for structured logging. If no active span is found in the context, it returns a Skip field.
//
// Parameters:
//   - ctx: The context containing the trace information
//
// Returns:
//   - zapcore.Field: A zap field containing the trace and span IDs
func Format(ctx context.Context) zapcore.Field {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return zap.Skip()
	}

	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()

	return zap.Inline(&traceLog{traceID, spanID})
}
