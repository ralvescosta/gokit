// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

import (
	"crypto/tls"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	// RMQConnection defines the interface for a RabbitMQ connection.
	// It abstracts the underlying AMQP connection and provides methods
	// for creating channels, retrieving connection state, and closing the connection.
	RMQConnection interface {
		// Channel creates a new channel on the connection.
		// Returns a pointer to an AMQP channel and any error encountered.
		Channel() (*amqp.Channel, error)

		// ConnectionState returns the TLS connection state if TLS is enabled.
		ConnectionState() tls.ConnectionState

		// Close gracefully closes the connection and all its channels.
		// It waits for confirmation from the server.
		Close() error
	}
)
