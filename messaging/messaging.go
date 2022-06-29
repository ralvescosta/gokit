package messaging

import (
	"github.com/ralvescostati/pkgs/messaging/kafka"
	"github.com/ralvescostati/pkgs/messaging/rabbitmq"
)

type (
	IMessageBroker interface {
		rabbitmq.IRabbitMQMessaging
		kafka.IKafkaMessaging
	}
)
