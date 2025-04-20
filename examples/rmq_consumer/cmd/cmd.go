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

		ctn.Logger.Debug("Stating RabbitMQ Consumer...")

		channel, err := rabbitmq.NewChannel(ctn.Cfg)
		if err != nil {
			ctn.Logger.Error("could not start rabbitmq client", zap.Error(err))
			return err
		}

		topology, err := rabbitmq.
			NewTopology(ctn.Cfg).
			Exchange(rabbitmq.NewFanoutExchange(ExchangeName)).
			Queue(rabbitmq.NewQueue(QueueName).WithDQL().WithRetry(30*time.Second, 3)).
			QueueBinding(rabbitmq.NewQueueBinding().Queue(QueueName).Exchange(ExchangeName)).
			Channel(channel).
			Apply()
		if err != nil {
			ctn.Logger.Error("topology error", zap.Error(err))
		}

		dispatcher := rabbitmq.NewDispatcher(ctn.Cfg, channel, topology.GetQueuesDefinition())
		ctn.BasicConsumer.Install(dispatcher)
		dispatcher.ConsumeBlocking(ctn.Sig)

		return nil
	},
}
