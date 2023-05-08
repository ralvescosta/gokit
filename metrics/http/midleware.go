package http

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
)

type (
	HTTPMetricsMiddleware interface {
		Handler(next http.Handler) http.Handler
	}

	httpMetricsMiddleware struct {
		meter           metric.Meter
		requestCounter  instrument.Int64Counter
		requestDuration instrument.Float64Histogram
	}

	responseWriter struct {
		http.ResponseWriter
		statusCode int
	}
)

func NewHTTPMetricsMiddleware() (HTTPMetricsMiddleware, error) {
	meter := global.Meter("github.com/ralvescosta/gokit/metric/http")

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

func (lrw *responseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
