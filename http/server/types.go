package server

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
)

type (
	HttpServerBuilder interface {
		WithTLS() HttpServerBuilder
		Timeouts(read, write, idle time.Duration) HttpServerBuilder
		WithProfiling() HttpServerBuilder

		Build()
		RegisterRoute(method string, path string, handler http.HandlerFunc) error
		Run() error
	}

	HttpServer struct {
		cfg           *env.Configs
		logger        logging.ILogger
		router        *chi.Mux
		server        *http.Server
		readTimeout   time.Duration
		writeTimeout  time.Duration
		idleTimeout   time.Duration
		withTLS       bool
		withProfiling bool
		sig           chan os.Signal
	}
)
