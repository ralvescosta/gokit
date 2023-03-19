package rabbitmq

import "github.com/streadway/amqp"

type (
	AMQPConnection interface {
		Channel() (*amqp.Channel, error)
	}
)
