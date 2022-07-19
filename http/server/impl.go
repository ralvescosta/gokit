package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func New(
	cfg *env.Configs,
	logger logging.ILogger,
	sig chan os.Signal,
) HttpServerBuilder {
	return &HttpServer{
		cfg:          cfg,
		logger:       logger,
		readTimeout:  5 * time.Second,
		writeTimeout: 10 * time.Second,
		idleTimeout:  30 * time.Second,
		sig:          sig,
	}
}

func (s *HttpServer) WithTLS() HttpServerBuilder {
	s.withTLS = true
	return s
}

func (s *HttpServer) Timeouts(read, write, idle time.Duration) HttpServerBuilder {
	s.readTimeout = read
	s.writeTimeout = write
	s.idleTimeout = idle
	return s
}

func (s *HttpServer) WithProfiling() HttpServerBuilder {
	s.withProfiling = true
	return s
}

func (s *HttpServer) Build() {
	s.router = chi.NewRouter()

	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Heartbeat("/heartbeat"))
	s.router.Use(middleware.AllowContentType("application/json"))
}

func (*HttpServer) RegisterRoute(method string, path string, handler http.HandlerFunc) error {
	return nil
}

func (s *HttpServer) Run() error {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	handler := otelhttp.NewHandler(s.router, "")

	s.server = &http.Server{
		Addr:         addr,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
		Handler:      handler,
	}

	return s.server.ListenAndServe()
}
