// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package cmd

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	configsBuilder "github.com/ralvescosta/gokit/configs_builder"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/mqtt"

	"github.com/ralvescosta/gokit/examples/mqtt_subscriber/pkg/consumers"
)

const (
	Topic = "my/topic/#"
)

type Container struct {
	Cfg           *configs.Configs
	Logger        logging.Logger
	Sig           chan os.Signal
	MQTTClient    mqtt.MQTTClient
	Dispatcher    mqtt.Dispatcher
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

	mqttClient := mqtt.NewMQTTClient(cfgs)
	dispatcher := mqtt.NewDispatcher(cfgs.Logger, mqttClient.Client())
	basicConsumer := consumers.NewBasicMessage(cfgs.Logger, Topic)

	return &Container{
		Cfg:           cfgs,
		Logger:        cfgs.Logger,
		Sig:           make(chan os.Signal, 1),
		MQTTClient:    mqttClient,
		Dispatcher:    dispatcher,
		BasicConsumer: basicConsumer,
	}, nil
}
