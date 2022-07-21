package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"syscall"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/http/server"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/telemetry/trace"
)

func main() {
	ctx := context.Background()

	cfg, err := env.
		New().
		Tracing().
		HTTPServer().
		Build()
	if err != nil {
		panic(err)
	}

	logger, err := logging.NewDefaultLogger(cfg)
	if err != nil {
		panic(err)
	}

	close, err := trace.
		NewOTLP(cfg, logger).
		WithApiKeyHeader().
		Build(ctx)
	if err != nil {
		panic(err)
	}
	defer close(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	httpServer := server.
		New(cfg, logger, sig).
		WithTracing().
		Build()

	httpServer.RegisterRoute(http.MethodGet, "/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("oi"))
	})

	httpServer.RegisterRoute(http.MethodGet, "/world", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("oi"))
	})

	if err := httpServer.Run(); err != nil {
		logger.Error("server error", logging.ErrorField(err))
	}
}
