// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package cmd

import (
	"time"

	"github.com/ralvescosta/gokit/rabbitmq"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var ConsumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consumer Command",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctn, err := NewContainer()
		if err != nil {
			return err
		}

		ctn.Logger.Debug("Starting RabbitMQ Consumer...")

		topology, err := rabbitmq.
			NewTopology(ctn.Cfg).
			Exchange(rabbitmq.NewFanoutExchange(ExchangeName)).
			Queue(rabbitmq.NewQueue(QueueName).WithDQL().WithRetry(30*time.Second, 3)).
			QueueBinding(rabbitmq.NewQueueBinding().Queue(QueueName).Exchange(ExchangeName)).
			Channel(ctn.AMQPChannel).
			Apply()
		if err != nil {
			ctn.Logger.Error("topology error", zap.Error(err))
		}

		dispatcher := rabbitmq.NewDispatcher(ctn.Cfg, ctn.AMQPChannel, topology.GetQueuesDefinition())
		ctn.BasicConsumer.Install(dispatcher)
		dispatcher.ConsumeBlocking(ctn.Sig)

		return nil
	},
}
