// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package tiny_http provides a lightweight, easy-to-use HTTP server implementation
// built on top of chi router with sensible defaults and a fluent interface for configuration.
package tiny_http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

type (
	// TinyServer defines the interface for a lightweight HTTP server with a fluent API
	// for configuring routes, middlewares, and other server features.
	TinyServer interface {
		// Sig sets a signal channel for graceful shutdown
		Sig(sig chan os.Signal) TinyServer

		// Prometheus enables the /metrics endpoint for Prometheus metrics
		Prometheus() TinyServer

		// Route registers a new route with the specified HTTP method, path, and handler
		Route(method string, path string, handler http.HandlerFunc) TinyServer

		// Middleware adds one or more middleware functions to the middleware stack
		Middleware(middlewares ...func(http.Handler) http.Handler) TinyServer

		// Run starts the HTTP server and blocks until the server is shutdown
		Run() error
	}

	// tinyHTTPServer is the implementation of the TinyServer interface
	tinyHTTPServer struct {
		cfg    *configs.HTTPConfigs
		logger logging.Logger
		router *chi.Mux
		server *http.Server
		sig    chan os.Signal
	}
)

// NewTinyServer creates a new TinyServer instance with sensible defaults
// It automatically configures common middleware like request ID, real IP,
// panic recovery, logging, content type validation, compression, and health check
func NewTinyServer(cfgs *configs.Configs) TinyServer {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.DefaultLogger)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.Compress(3, "application/json"))
	router.Use(middleware.Heartbeat("/health"))

	return &tinyHTTPServer{
		cfg:    cfgs.HTTPConfigs,
		logger: cfgs.Logger,
		router: router,
	}
}

// Sig sets a signal channel for graceful shutdown
// The server will listen for signals on this channel and initiate a graceful
// shutdown when a signal is received
func (t *tinyHTTPServer) Sig(sig chan os.Signal) TinyServer {
	t.sig = sig
	return t
}

// Prometheus enables the /metrics endpoint for Prometheus metrics
// This endpoint will expose runtime and application metrics in Prometheus format
func (t *tinyHTTPServer) Prometheus() TinyServer {
	t.router.Method(http.MethodGet, "/metrics", promhttp.Handler())
	return t
}

// Route registers a new route with the specified HTTP method, path, and handler
// This method is used to define API endpoints
func (t *tinyHTTPServer) Route(method string, path string, handler http.HandlerFunc) TinyServer {
	t.router.Method(method, path, handler)
	return t
}

// Middleware adds one or more middleware functions to the middleware stack
// Middlewares are executed in the order they are added
func (t *tinyHTTPServer) Middleware(middlewares ...func(http.Handler) http.Handler) TinyServer {
	t.router.Use(middlewares...)
	return t
}

// Run starts the HTTP server and blocks until the server is shutdown
// It sets up graceful shutdown and handles server errors
func (t *tinyHTTPServer) Run() error {
	t.logger.Debug("[gokit:tiny_server] starting http server...")
	t.server = &http.Server{
		Addr:         t.cfg.Addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      t.router,
	}

	t.logger.Debug("[gokit:tiny_server] configuring graceful shutdown...")
	ctx, ctxCancelFunc := context.WithCancel(context.Background())
	go t.shutdown(ctx, ctxCancelFunc)

	t.logger.Info(fmt.Sprintf("[gokit:tiny_server] %s started", t.cfg.Addr))
	if err := t.server.ListenAndServe(); err != nil {
		t.logger.Error("[gokit:tiny_server] http server error", zap.Error(err))
		return err
	}

	<-ctx.Done()

	return nil
}

// shutdown handles the graceful shutdown process when a signal is received
// It attempts to close existing connections gracefully and has a timeout
// for shutting down
func (t *tinyHTTPServer) shutdown(ctx context.Context, ctxCancelFunc context.CancelFunc) {
	if t.sig == nil {
		return
	}

	<-t.sig

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			t.logger.Fatal("[gokit:tiny_server] graceful shutdown timed out.. forcing exit.")
		}
	}()

	err := t.server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatal(err)
	}

	ctxCancelFunc()
}
