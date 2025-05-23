// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package server provides components for building and managing HTTP servers.
// It offers a flexible API for defining routes, registering middleware, and
// creating server groups with shared configurations.
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
	"github.com/ralvescosta/gokit/logging"
	metrics "github.com/ralvescosta/gokit/metrics/http"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"

	"github.com/ralvescosta/gokit/httpw"
)

type (
	// HTTPServer defines the interface for managing an HTTP server.
	// It provides methods for registering routes, middleware, and
	// starting the server with graceful shutdown capabilities.
	HTTPServer interface {
		// BasicRoute registers a basic HTTP route with the specified method, path, and handler.
		BasicRoute(method string, path string, handler http.HandlerFunc) error

		// Route registers a Route object with the server.
		Route(r *Route) error

		// Middleware adds middleware to the server for all routes.
		Middleware(m *Middleware)

		// Group creates a routing group with a shared URL prefix and route definitions.
		Group(pattern string, routes []*Route) error

		// Run starts the HTTP server and returns an error if the server fails to start.
		Run() error
	}

	// httpServer implements the HTTPServer interface.
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

// Map of allowed HTTP methods for route registration
var (
	allowedMethod = map[string]uint8{
		http.MethodPost: 1, http.MethodGet: 1, http.MethodPatch: 1, http.MethodPut: 1, http.MethodDelete: 1,
	}
)

// BasicRoute registers a basic HTTP route with the specified method, path, and handler.
func (s *httpServer) BasicRoute(method string, path string, handler http.HandlerFunc) error {
	if method != "" {
		return s.registerRoute(s.router, NewRouteBuilder().Method(method).Path(path).Handler(handler).Build(), "")
	}

	return httpw.ErrHTTPMethodMethodIsRequired
}

// Route registers a Route object with the server.
func (s *httpServer) Route(r *Route) error {
	if r.method != "" {
		return s.registerRoute(s.router, r, "")
	}

	return httpw.ErrHTTPMethodMethodIsRequired
}

// Group creates a routing group with a shared URL prefix and route definitions.
func (s *httpServer) Group(group string, routes []*Route) error {
	var err error

	s.router.Route(group, func(router chi.Router) {
		for _, route := range routes {
			if route.method == "" {
				err = httpw.ErrHTTPMethodMethodIsRequired
				break
			}

			err = s.registerRoute(router, route, group)
			if err != nil {
				break
			}
		}
	})

	return err
}

// Middleware adds middleware to the server for all routes.
func (s *httpServer) Middleware(m *Middleware) {
	s.router.Use(m.middlewares...)
}

// Run starts the HTTP server and listens for incoming requests.
// It also sets up graceful shutdown through a signal channel.
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

// registerRoute registers an HTTP route with the router and applies any middleware or tracing.
// It validates the HTTP method and returns an error if the method is not supported.
func (s *httpServer) registerRoute(router chi.Router, route *Route, group string) error {
	if _, ok := allowedMethod[route.method]; !ok {
		s.logger.Warn(httpw.Message("method not allowed"))
		return httpw.ErrInvalidHTTPMethod
	}

	s.logger.Debug(s.logRouterRegister(route.method, fmt.Sprintf("%v%v", group, route.path)))
	var newHandler http.Handler = route.handler
	if s.withTracing {
		newHandler = otelhttp.NewHandler(route.handler, otlpOperationName(route.method, route.path))
	}

	if len(route.middlewares) > 0 {
		router.With(route.middlewares...).Method(route.method, route.path, newHandler)
		return nil
	}
	router.Method(route.method, route.path, newHandler)

	s.logger.Debug(httpw.Message("router registered"))
	return nil
}

// shutdown handles graceful shutdown of the server when a signal is received.
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

// logRouterRegister formats a log message for route registration.
func (s *httpServer) logRouterRegister(method, path string) string {
	return httpw.Message(fmt.Sprintf("registering route: %s %s", method, path))
}

// otlpOperationName formats an OpenTelemetry operation name for a route.
func otlpOperationName(method, path string) string {
	return method + " " + path
}
