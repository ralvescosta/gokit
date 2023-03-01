package httpw

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
)

type (
	HTTPServerBuilder interface {
		WithTLS() HTTPServerBuilder
		Timeouts(read, write, idle time.Duration) HTTPServerBuilder
		WithProfiling() HTTPServerBuilder
		WithTracing() HTTPServerBuilder
		WithMetrics(metricKind MetricKind) HTTPServerBuilder
		Build() HTTPServer
	}

	HTTPServer interface {
		RegisterRoute(method string, path string, handler http.HandlerFunc) error
		RegisterPrometheus()
		Run() error
	}

	MetricKind  int
	TracingKind int

	httpServerImpl struct {
		cfg           *env.HTTPConfigs
		logger        logging.Logger
		router        *chi.Mux
		server        *http.Server
		sig           chan os.Signal
		readTimeout   time.Duration
		writeTimeout  time.Duration
		idleTimeout   time.Duration
		withTLS       bool
		withProfiling bool
		withTracing   bool
		withMetric    bool
		metricKind    MetricKind
	}
)
