// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package http provides HTTP middleware for metrics collection and monitoring.
package http

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type (
	// HTTPMetricsMiddleware defines an interface for HTTP metrics collection middleware.
	// It provides a way to collect and report metrics for HTTP requests.
	HTTPMetricsMiddleware interface {
		// Handler wraps an existing http.Handler with metrics collection.
		// It tracks request counts and durations with attributes for method, URI, and status code.
		Handler(next http.Handler) http.Handler
	}

	// httpMetricsMiddleware implements the HTTPMetricsMiddleware interface.
	httpMetricsMiddleware struct {
		// meter is the OpenTelemetry meter used to create metrics instruments.
		meter metric.Meter

		// requestCounter counts the number of HTTP requests.
		requestCounter metric.Int64Counter

		// requestDuration measures the duration of HTTP requests.
		requestDuration metric.Float64Histogram
	}

	// responseWriter wraps an http.ResponseWriter to capture the status code.
	responseWriter struct {
		http.ResponseWriter
		statusCode int
	}
)

// NewHTTPMetricsMiddleware creates a new HTTP metrics middleware that collects
// request counts and durations for HTTP requests.
//
// Returns:
//   - An HTTPMetricsMiddleware interface for HTTP metrics collection.
//   - An error if the meter instruments cannot be created.
func NewHTTPMetricsMiddleware() (HTTPMetricsMiddleware, error) {
	meter := otel.Meter("github.com/ralvescosta/gokit/metric/http")

	counter, err := meter.Int64Counter("http.requests", metric.WithDescription("HTTP Requests Counter"))
	if err != nil {
		return nil, err
	}

	duration, err := meter.Float64Histogram("http.request.duration", metric.WithDescription("HTTP Request Duration"))
	if err != nil {
		return nil, err
	}

	return &httpMetricsMiddleware{
		meter:           meter,
		requestCounter:  counter,
		requestDuration: duration,
	}, nil
}

// Handler wraps an HTTP handler with metrics collection.
// It records the request duration and increments the request counter
// with method, URI, and status code attributes.
//
// Parameters:
//   - next: The HTTP handler to wrap with metrics collection.
//
// Returns:
//   - An HTTP handler that collects metrics before calling the wrapped handler.
func (m *httpMetricsMiddleware) Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rw := &responseWriter{w, http.StatusOK}
		start := time.Now()

		next.ServeHTTP(rw, r.WithContext(ctx))

		m.requestDuration.Record(
			ctx,
			float64(time.Since(start).Nanoseconds()),
			metric.WithAttributes(
				attribute.String("method", r.Method),
				attribute.String("uri", r.RequestURI),
				attribute.Int("statusCode", rw.statusCode),
			),
		)

		m.requestCounter.Add(
			ctx,
			1,
			metric.WithAttributes(
				attribute.String("method", r.Method),
				attribute.String("uri", r.RequestURI),
				attribute.Int("statusCode", rw.statusCode),
			),
		)
	}

	return http.HandlerFunc(fn)
}

// WriteHeader captures the status code and delegates to the wrapped ResponseWriter.
//
// Parameters:
//   - code: The HTTP status code to write.
func (lrw *responseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
