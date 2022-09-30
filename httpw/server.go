package httpw

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewServer(
	cfg *env.Config,
	logger logging.Logger,
	sig chan os.Signal,
) HTTPServerBuilder {
	return &HTTPServerImpl{
		cfg:          cfg,
		logger:       logger,
		readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,
		idleTimeout:  30 * time.Second,
		sig:          sig,
	}
}

func (s *HTTPServerImpl) WithTLS() HTTPServerBuilder {
	s.withTLS = true
	return s
}

func (s *HTTPServerImpl) Timeouts(read, write, idle time.Duration) HTTPServerBuilder {
	s.readTimeout = read
	s.writeTimeout = write
	s.idleTimeout = idle
	return s
}

func (s *HTTPServerImpl) WithProfiling() HTTPServerBuilder {
	s.withProfiling = true
	return s
}

func (s *HTTPServerImpl) WithTracing() HTTPServerBuilder {
	s.withTracing = true
	return s
}

func (s *HTTPServerImpl) WithMetrics(kind MetricKind) HTTPServerBuilder {
	s.withMetric = true
	s.metricKind = kind

	return s
}

func (s *HTTPServerImpl) Build() HTTPServer {
	s.logger.Debug(Message("creating the server..."))
	s.router = chi.NewRouter()

	if s.withMetric {
		s.registerMetricMiddleware()
	}

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.DefaultLogger)
	s.router.Use(middleware.Heartbeat("/health"))
	s.router.Use(middleware.AllowContentType("application/json"))

	if s.withProfiling {
		s.router.Mount("/debug", middleware.Profiler())
	}

	if s.withMetric {
		s.installMetrics()
	}

	s.logger.Debug(Message("server was created"))
	return s
}

func (s *HTTPServerImpl) RegisterRoute(method string, path string, handler http.HandlerFunc) error {
	if _, ok := allowedHTTPMethods[method]; !ok {
		s.logger.Warn(Message("method not allowed"))
		return ErrorInvalidHttpMethod
	}

	s.logger.Debug(LogRouterRegister(method, path))
	var newHandler http.Handler = handler
	if s.withTracing {
		newHandler = otelhttp.NewHandler(handler, OTLPOperationName(method, path))
	}

	s.router.Method(method, path, newHandler)

	s.logger.Debug(Message("router registered"))
	return nil
}

func (s *HTTPServerImpl) Run() error {
	s.logger.Debug(Message("starting http server..."))

	s.server = &http.Server{
		Addr:         s.cfg.HTTP_ADDR,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
		Handler:      s.router,
	}

	s.logger.Debug(Message("configuring graceful shutdown..."))
	ctx, ctxCancelFunc := context.WithCancel(context.Background())
	go s.shutdown(ctx, ctxCancelFunc)

	s.logger.Info(Message(fmt.Sprintf("%s started", s.cfg.HTTP_ADDR)))
	if err := s.server.ListenAndServe(); err != nil {
		s.logger.Error(Message("http server error"), logging.ErrorField(err))
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *HTTPServerImpl) registerMetricMiddleware() {
	switch s.metricKind {
	case PrometheusMetricKind:
		s.router.Use(PrometheusMiddleware)
	case OtelMetricKind:
		s.logger.Warn("WITHOUT METRICS")
	}
}

func (s *HTTPServerImpl) installMetrics() {
	s.logger.Debug(Message("Installing metrics..."))

	switch s.metricKind {
	case PrometheusMetricKind:
		s.installPrometheus()
	case OtelMetricKind:
		s.logger.Info("otel is not implemented yet")
	}

	s.logger.Debug(Message("metrics installed"))
}

func (s *HTTPServerImpl) installPrometheus() {
	handler := promhttp.Handler()
	method := http.MethodGet
	pattern := "/metrics"

	if s.withTracing {
		handler = otelhttp.NewHandler(promhttp.Handler(), OTLPOperationName(method, pattern))
	}

	s.router.Method(method, pattern, handler)
}

func (s *HTTPServerImpl) shutdown(ctx context.Context, ctxCancelFunc context.CancelFunc) {
	<-s.sig

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			s.logger.Fatal("graceful shutdown timed out.. forcing exit.")
		}
	}()

	err := s.server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatal(err)
	}

	ctxCancelFunc()
}
