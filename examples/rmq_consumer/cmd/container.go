// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package cmd

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	configsBuilder "github.com/ralvescosta/gokit/configs_builder"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/rabbitmq"
	"go.uber.org/zap"

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
	AMQPChannel   rabbitmq.AMQPChannel
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

	channel, err := rabbitmq.NewChannel(cfgs)
	if err != nil {
		cfgs.Logger.Error("could not start rabbitmq client", zap.Error(err))
		return nil, err
	}

	basicConsumer := consumers.NewBasicMessage(cfgs.Logger, QueueName)

	return &Container{
		Cfg:           cfgs,
		Logger:        cfgs.Logger,
		Sig:           make(chan os.Signal, 1),
		AMQPChannel:   channel,
		BasicConsumer: basicConsumer,
	}, nil
}
