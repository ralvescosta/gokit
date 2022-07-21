package server

import (
	"context"
	"log"
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

	if s.withProfiling {
		s.router.Mount("/debug", middleware.Profiler())
	}
}

func (s *HttpServer) RegisterRoute(method string, path string, handler http.HandlerFunc) error {
	switch method {
	case http.MethodGet:
		s.router.Get(path, handler)
	case http.MethodPost:
		s.router.Post(path, handler)
	case http.MethodPut:
		s.router.Put(path, handler)
	case http.MethodPatch:
		s.router.Patch(path, handler)
	case http.MethodDelete:
		s.router.Delete(path, handler)
	default:
		return ErrorInvalidHttpMethod
	}

	return nil
}

func (s *HttpServer) Run() error {
	handler := otelhttp.NewHandler(s.router, "")

	s.server = &http.Server{
		Addr:         s.cfg.HTTP_ADDR,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
		Handler:      handler,
	}

	ctx, ctxCancelFunc := context.WithCancel(context.Background())

	go s.shutdown(ctx, ctxCancelFunc)

	if err := s.server.ListenAndServe(); err != nil {
		s.logger.Error(LogMessage("Listen and Serve"), logging.ErrorField(err))
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *HttpServer) shutdown(ctx context.Context, ctxCancelFunc context.CancelFunc) {
	<-s.sig

	shutdownCtx, _ := context.WithTimeout(ctx, 30*time.Second)
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
