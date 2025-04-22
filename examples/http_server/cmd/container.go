// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ralvescosta/gokit/auth/auth0"
	"github.com/ralvescosta/gokit/configs"
	configsBuilder "github.com/ralvescosta/gokit/configs_builder"
	"github.com/ralvescosta/gokit/httpw/middlewares"
	"github.com/ralvescosta/gokit/httpw/server"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/metrics"
	"github.com/ralvescosta/gokit/tracing"

	"github.com/ralvescosta/gokit/examples/http_server/internals/repositories"
	"github.com/ralvescosta/gokit/examples/http_server/internals/services"
	"github.com/ralvescosta/gokit/examples/http_server/pkg/handlers"
)

type Container struct {
	Cfg           *configs.Configs
	Logger        logging.Logger
	Sig           chan os.Signal
	HTTPServer    server.HTTPServer
	BooksHandlers handlers.HTTPHandlers
}

func NewContainer() (*Container, error) {
	cfgs, err := configsBuilder.
		NewConfigsBuilder().
		HTTP().
		Tracing().
		Metrics().
		Build()

	if err != nil {
		return nil, err
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	httpServer := server.
		NewHTTPServerBuilder(cfgs).
		Signal(signalChan).
		WithMetrics().
		WithTracing().
		WithOpenAPI().
		Build()

	identity := auth0.NewAuth0TokenManger(cfgs)
	authMiddleware := middlewares.NewAuthorization(cfgs.Logger, identity)

	booksRepository := repositories.NewBookRepository(cfgs)
	booksSvc := services.NewBookService(booksRepository)
	booksHandlers := handlers.NewHandler(cfgs.Logger, booksSvc, authMiddleware)

	tracing.
		NewOTLP(cfgs).
		Build()

	metrics.
		NewOTLPBuilder(cfgs)

	return &Container{
		Cfg:           cfgs,
		Logger:        cfgs.Logger,
		Sig:           signalChan,
		HTTPServer:    httpServer,
		BooksHandlers: booksHandlers,
	}, nil
}
