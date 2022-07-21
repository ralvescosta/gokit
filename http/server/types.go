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
	HTTPServerBuilder interface {
		WithTLS() HTTPServerBuilder
		Timeouts(read, write, idle time.Duration) HTTPServerBuilder
		WithProfiling() HTTPServerBuilder
		WithTracing() HTTPServerBuilder
		Build() IHTTPServer
	}

	IHTTPServer interface {
		RegisterRoute(method string, path string, handler http.HandlerFunc) error
		Run() error
	}

	HTTPServer struct {
		cfg           *env.Configs
		logger        logging.ILogger
		router        *chi.Mux
		server        *http.Server
		readTimeout   time.Duration
		writeTimeout  time.Duration
		idleTimeout   time.Duration
		withTLS       bool
		withProfiling bool
		withTracing   bool
		sig           chan os.Signal
	}
)
