package tinyhttp

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
	TinyServer interface {
		Sig(sig chan os.Signal) TinyServer
		Prometheus() TinyServer
		Route(method string, path string, handler http.HandlerFunc) TinyServer
		Middleware(middlewares ...func(http.Handler) http.Handler) TinyServer
		Run() error
	}

	tinyHTTPServer struct {
		cfg    *configs.HTTPConfigs
		logger logging.Logger
		router *chi.Mux
		server *http.Server
		sig    chan os.Signal
	}
)

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

func (t *tinyHTTPServer) Sig(sig chan os.Signal) TinyServer {
	t.sig = sig
	return t
}

func (t *tinyHTTPServer) Prometheus() TinyServer {
	t.router.Method(http.MethodGet, "/metrics", promhttp.Handler())
	return t
}

func (t *tinyHTTPServer) Route(method string, path string, handler http.HandlerFunc) TinyServer {
	t.router.Method(method, path, handler)
	return t
}

func (t *tinyHTTPServer) Middleware(middlewares ...func(http.Handler) http.Handler) TinyServer {
	t.router.Use(middlewares...)
	return t
}

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
