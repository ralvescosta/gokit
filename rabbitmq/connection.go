package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type (
	AMQPConnection interface {
		Channel() (*amqp.Channel, error)
	}
)
