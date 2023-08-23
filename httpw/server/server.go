package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/httpw"
	"github.com/ralvescosta/gokit/logging"
	metrics "github.com/ralvescosta/gokit/metrics/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

type (
	HTTPServer interface {
		BasicRoute(method string, path string, handler http.HandlerFunc) error
		Route(r *Route) error
		Group(pattern string, routes []*Route) error
		Run() error
	}

	httpServer struct {
		cfg               *configs.HTTPConfigs
		logger            logging.Logger
		router            *chi.Mux
		server            *http.Server
		sig               chan os.Signal
		metricsMiddleware metrics.HTTPMetricsMiddleware
		readTimeout       time.Duration
		writeTimeout      time.Duration
		idleTimeout       time.Duration
		withTracing       bool
	}
)

var (
	allowedMethod = map[string]uint8{
		http.MethodPost: 1, http.MethodGet: 1, http.MethodPatch: 1, http.MethodPut: 1, http.MethodDelete: 1,
	}
)

func (s *httpServer) BasicRoute(method string, path string, handler http.HandlerFunc) error {
	return s.registerRoute(s.router, method, "", path, handler)
}

func (s *httpServer) Route(r *Route) error {
	if r.middlewares != nil {
		s.router.Use(r.middlewares...)
	}

	if r.method != "" {
		return s.registerRoute(s.router, r.method, "", r.path, r.handler)
	}

	return nil
}

func (s *httpServer) Group(pattern string, routes []*Route) error {
	var err error

	s.router.Route(pattern, func(r chi.Router) {
		for _, v := range routes {
			if v.method != "" {
				err = s.registerRoute(r, v.method, pattern, v.path, v.handler)
			} else {
				r.Use(v.middlewares...)
			}

			if err != nil {
				break
			}
		}
	})

	return err
}

func (s *httpServer) Run() error {
	s.logger.Debug(httpw.Message("starting http server..."))

	s.server = &http.Server{
		Addr:         s.cfg.Addr,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
		Handler:      s.router,
	}

	s.logger.Debug(httpw.Message("configuring graceful shutdown..."))
	ctx, ctxCancelFunc := context.WithCancel(context.Background())
	go s.shutdown(ctx, ctxCancelFunc)

	s.logger.Info(httpw.Message(fmt.Sprintf("%s started", s.cfg.Addr)))
	if err := s.server.ListenAndServe(); err != nil {
		s.logger.Error(httpw.Message("http server error"), zap.Error(err))
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *httpServer) registerRoute(r chi.Router, method, pattern, path string, handler http.HandlerFunc) error {
	if _, ok := allowedMethod[method]; !ok {
		s.logger.Warn(httpw.Message("method not allowed"))
		return httpw.InvalidHttpMethodError
	}

	s.logger.Debug(s.logRouterRegister(method, fmt.Sprintf("%v%v", pattern, path)))
	var newHandler http.Handler = handler
	if s.withTracing {
		newHandler = otelhttp.NewHandler(handler, otlpOperationName(method, path))
	}

	r.Method(method, path, newHandler)

	s.logger.Debug(httpw.Message("router registered"))
	return nil
}

func (s *httpServer) shutdown(ctx context.Context, ctxCancelFunc context.CancelFunc) {
	if s.sig == nil {
		return
	}

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

func (s *httpServer) logRouterRegister(method, path string) string {
	return httpw.Message(fmt.Sprintf("registering route: %s %s", method, path))
}

func otlpOperationName(method, path string) string {
	return method + " " + path
}
