package tinyHttp

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
)

type (
	TinyServer interface {
		Sig(sig chan os.Signal) TinyServer
		Prometheus() TinyServer
		Route(method string, path string, handler http.HandlerFunc) TinyServer
		Middleware(middlewares ...func(http.Handler) http.Handler) TinyServer
	}

	tinyHTTPServer struct {
		cfg    *configs.HTTPConfigs
		logger logging.Logger
		router *chi.Mux
		sig    chan os.Signal
	}
)

func NewTinyServer(cfg *configs.HTTPConfigs, logger logging.Logger) TinyServer {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.DefaultLogger)
	router.Use(middleware.AllowContentType("application/json"))
	router.Use(middleware.Compress(3, "application/json"))
	router.Use(middleware.Heartbeat("/health"))

	return &tinyHTTPServer{
		cfg:    cfg,
		logger: logger,
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
