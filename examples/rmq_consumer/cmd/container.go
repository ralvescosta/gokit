// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package cmd

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	configsBuilder "github.com/ralvescosta/gokit/configs_builder"
	"github.com/ralvescosta/gokit/logging"

	"github.com/ralvescosta/gokit/examples/rmq_consumer/pkg/consumers"
)

const (
	QueueName    = "observability.queue"
	ExchangeName = "observability"
)

type Container struct {
	Cfg           *configs.Configs
	Logger        logging.Logger
	Sig           chan os.Signal
	BasicConsumer consumers.BasicConsumer
}

func NewContainer() (*Container, error) {
	cfgs, err := configsBuilder.
		NewConfigsBuilder().
		RabbitMQ().
		Tracing().
		Metrics().
		Build()

	if err != nil {
		return nil, err
	}

	basicConsumer := consumers.NewBasicMessage(cfgs.Logger, QueueName)

	return &Container{
		Cfg:           cfgs,
		Logger:        cfgs.Logger,
		Sig:           make(chan os.Signal, 1),
		BasicConsumer: basicConsumer,
	}, nil
}
