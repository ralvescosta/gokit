package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/httpw"
	"github.com/ralvescosta/gokit/logging"
	metrics "github.com/ralvescosta/gokit/metrics/http"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type (
	MetricKind  int
	TracingKind int

	HTTPServerBuilder interface {
		WithTLS() HTTPServerBuilder
		Timeouts(read, write, idle time.Duration) HTTPServerBuilder
		WithTracing() HTTPServerBuilder
		WithMetrics() HTTPServerBuilder
		//Doc will be available only in local, develop and staging environment
		WithOpenAPI() HTTPServerBuilder
		Signal(sig chan os.Signal) HTTPServerBuilder
		ExportPrometheusScraping() HTTPServerBuilder
		Build() HTTPServer
	}

	httpServerBuilder struct {
		env                      configs.Environment
		cfg                      *configs.HTTPConfigs
		logger                   logging.Logger
		sig                      chan os.Signal
		readTimeout              time.Duration
		writeTimeout             time.Duration
		idleTimeout              time.Duration
		withTLS                  bool
		withTracing              bool
		withMetric               bool
		exportPrometheusScraping bool
		withOpenApi              bool
		metricKind               MetricKind
	}
)

func NewHTTPServerBuilder(cfg *configs.HTTPConfigs, logger logging.Logger) HTTPServerBuilder {
	return &httpServerBuilder{
		cfg:          cfg,
		logger:       logger,
		readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,
		idleTimeout:  30 * time.Second,
	}
}

func (s *httpServerBuilder) WithTLS() HTTPServerBuilder {
	s.withTLS = true
	return s
}

func (s *httpServerBuilder) Timeouts(read, write, idle time.Duration) HTTPServerBuilder {
	s.readTimeout = read
	s.writeTimeout = write
	s.idleTimeout = idle
	return s
}

func (s *httpServerBuilder) WithTracing() HTTPServerBuilder {
	s.withTracing = true
	return s
}

func (s *httpServerBuilder) WithMetrics() HTTPServerBuilder {
	s.withMetric = true

	return s
}

func (s *httpServerBuilder) WithOpenAPI() HTTPServerBuilder {
	s.withOpenApi = true
	return s
}

func (s *httpServerBuilder) ExportPrometheusScraping() HTTPServerBuilder {
	s.exportPrometheusScraping = true
	return s
}

func (s *httpServerBuilder) Signal(sig chan os.Signal) HTTPServerBuilder {
	s.sig = sig
	return s
}

func (s *httpServerBuilder) Build() HTTPServer {
	s.logger.Debug(httpw.Message("creating the server..."))

	server := httpServer{
		router:       chi.NewRouter(),
		logger:       s.logger,
		cfg:          s.cfg,
		sig:          s.sig,
		withTracing:  s.withMetric,
		readTimeout:  s.readTimeout,
		writeTimeout: s.writeTimeout,
		idleTimeout:  s.idleTimeout,
	}

	if s.withMetric {
		metricsMiddleware, _ := metrics.NewHTTPMetricsMiddleware()
		server.metricsMiddleware = metricsMiddleware
		server.router.Use(server.metricsMiddleware.Handler)
	}

	server.router.Use(middleware.RequestID)
	server.router.Use(middleware.RealIP)
	server.router.Use(middleware.Recoverer)
	server.router.Use(middleware.DefaultLogger)
	server.router.Use(middleware.Heartbeat("/health"))
	server.router.Use(middleware.AllowContentType("application/json"))

	if s.cfg.EnableProfiling {
		server.router.Mount("/debug", middleware.Profiler())
	}

	if s.exportPrometheusScraping {
		s.prometheusScrapingEndpoint(&server)
	}

	if s.withOpenApi {
		s.openAPIEndpoint(&server)
	}

	s.logger.Debug(httpw.Message("server was created"))
	return &server
}

func (sb *httpServerBuilder) prometheusScrapingEndpoint(s *httpServer) {
	handler := promhttp.Handler()
	method := http.MethodGet
	pattern := "/metrics"

	if sb.withTracing {
		handler = otelhttp.NewHandler(promhttp.Handler(), otlpOperationName(method, pattern))
	}

	s.router.Method(method, pattern, handler)
}

func (sb *httpServerBuilder) openAPIEndpoint(s *httpServer) {
	if sb.env == configs.PRODUCTION_ENV {
		return
	}

	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", sb.cfg.Addr)),
	))
}
