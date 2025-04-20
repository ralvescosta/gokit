// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type (
	AMQPConnection interface {
		Channel() (*amqp.Channel, error)
	}
)
