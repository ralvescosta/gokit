// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

import (
	"crypto/tls"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	RMQConnection interface {
		Channel() (*amqp.Channel, error)
		ConnectionState() tls.ConnectionState
		Close() error
	}
)
