package server

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

func New(
	cfg *env.Configs,
	logger logging.ILogger,
	sig chan os.Signal,
) HTTPServerBuilder {
	return &HTTPServer{
		cfg:          cfg,
		logger:       logger,
		readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,
		idleTimeout:  30 * time.Second,
		sig:          sig,
	}
}

func (s *HTTPServer) WithTLS() HTTPServerBuilder {
	s.withTLS = true
	return s
}

func (s *HTTPServer) Timeouts(read, write, idle time.Duration) HTTPServerBuilder {
	s.readTimeout = read
	s.writeTimeout = write
	s.idleTimeout = idle
	return s
}

func (s *HTTPServer) WithProfiling() HTTPServerBuilder {
	s.withProfiling = true
	return s
}

func (s *HTTPServer) WithTracing() HTTPServerBuilder {
	s.withTracing = true
	return s
}

func (s *HTTPServer) WithMetrics(kind MetricKind) HTTPServerBuilder {
	s.withMetric = true
	s.metricKind = kind

	return s
}

func (s *HTTPServer) Build() IHTTPServer {
	s.logger.Debug(LogMessage("creating the server..."))
	s.router = chi.NewRouter()

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.DefaultLogger)
	s.router.Use(middleware.Heartbeat("/heartbeat"))
	s.router.Use(middleware.AllowContentType("application/json"))

	if s.withProfiling {
		s.router.Mount("/debug", middleware.Profiler())
	}

	if s.withMetric {
		s.installMetrics()
	}

	s.logger.Debug(LogMessage("server was created"))
	return s
}

func (s *HTTPServer) RegisterRoute(method string, path string, handler http.HandlerFunc) error {
	if _, ok := allowedHTTPMethods[method]; !ok {
		s.logger.Warn(LogMessage("method not allowed"))
		return ErrorInvalidHttpMethod
	}

	s.logger.Debug(LogRouterRegister(method, path))
	var newHandler http.Handler = handler
	if s.withTracing {
		newHandler = otelhttp.NewHandler(handler, OTLPOperationName(method, path))
	}

	s.router.Method(method, path, newHandler)

	s.logger.Debug(LogMessage("router registered"))
	return nil
}

func (s *HTTPServer) Run() error {
	s.logger.Debug(LogMessage("starting http server..."))

	s.server = &http.Server{
		Addr:         s.cfg.HTTP_ADDR,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
		Handler:      s.router,
	}

	s.logger.Debug(LogMessage("configuring graceful shutdown..."))
	ctx, ctxCancelFunc := context.WithCancel(context.Background())
	go s.shutdown(ctx, ctxCancelFunc)

	s.logger.Info(LogMessage(fmt.Sprintf("%s started", s.cfg.HTTP_ADDR)))
	if err := s.server.ListenAndServe(); err != nil {
		s.logger.Error(LogMessage("http server error"), logging.ErrorField(err))
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *HTTPServer) installMetrics() {
	s.logger.Debug(LogMessage("Installing metrics..."))

	if s.metricKind != PrometheusMetricKind {
		s.logger.Error(LogMessage("MetricKind not allowed"))
		return
	}

	handler := promhttp.Handler()
	method := http.MethodGet
	pattern := "/metrics"

	if s.withTracing {
		handler = otelhttp.NewHandler(promhttp.Handler(), OTLPOperationName(method, pattern))
	}

	s.router.Method(method, pattern, handler)
}

func (s *HTTPServer) shutdown(ctx context.Context, ctxCancelFunc context.CancelFunc) {
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
