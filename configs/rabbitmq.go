// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// RabbitMQConfigs defines configuration parameters for RabbitMQ message broker connections.
// It contains all necessary connection details for establishing a connection to a RabbitMQ server.
type RabbitMQConfigs struct {
	// Schema defines the connection protocol (typically "amqp" or "amqps" for TLS)
	Schema string
	// Host specifies the RabbitMQ server hostname or IP address
	Host string
	// Port defines the network port on which the RabbitMQ server is listening
	Port string
	// User specifies the username for RabbitMQ server authentication
	User string
	// Password contains the authentication credential for the RabbitMQ user
	Password string
	// VHost specifies the virtual host to use on the RabbitMQ server
	VHost string
}
