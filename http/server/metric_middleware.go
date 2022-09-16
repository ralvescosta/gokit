package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *responseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

var duration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: "http",
	Name:      "request_duration_seconds",
	Help:      "The latency of the HTTP requests.",
	Buckets:   prometheus.DefBuckets,
}, []string{"handler", "method", "code"})

var amount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "http",
	Name:      "request_amount",
	Help:      "The amount of HTTP request.",
}, []string{"handler", "method", "code"})

func init() {
	prometheus.Register(duration)
	prometheus.Register(amount)
}

func MetricMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rw := NewResponseWriter(w)
		start := time.Now()

		next.ServeHTTP(rw, r.WithContext(ctx))

		duration.WithLabelValues(r.RequestURI, r.Method, strconv.Itoa(rw.statusCode)).Observe(float64(time.Since(start).Nanoseconds()))
		amount.WithLabelValues(r.RequestURI, r.Method, strconv.Itoa(rw.statusCode)).Inc()
	}

	return http.HandlerFunc(fn)
}
