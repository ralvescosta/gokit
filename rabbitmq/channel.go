// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ralvescosta/gokit/configs"
	"go.uber.org/zap"
)

type (
	// AMQPChannel defines the interface for a RabbitMQ channel.
	// It abstracts the operations that can be performed on a channel such as
	// declaring exchanges and queues, binding them, and publishing or consuming messages.
	AMQPChannel interface {
		// ExchangeDeclare declares an exchange on the channel.
		// The exchange will be created if it doesn't already exist.
		// Parameters:
		//   - name: The name of the exchange
		//   - kind: The exchange type (direct, fanout, topic, headers)
		//   - durable: Survive broker restarts
		//   - autoDelete: Delete when no longer used
		//   - internal: Can only be published to by other exchanges
		//   - noWait: Don't wait for a server confirmation
		//   - args: Additional arguments
		ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error

		// ExchangeBind binds an exchange to another exchange.
		// Parameters:
		//   - destination: The name of the destination exchange
		//   - key: The routing key to use
		//   - source: The name of the source exchange
		//   - noWait: Don't wait for a server confirmation
		//   - args: Additional arguments
		ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error

		// QueueDeclare declares a queue on the channel.
		// The queue will be created if it doesn't already exist.
		// Parameters:
		//   - name: The name of the queue
		//   - durable: Survive broker restarts
		//   - autoDelete: Delete when no longer used
		//   - exclusive: Used by only one connection and deleted when that connection closes
		//   - noWait: Don't wait for a server confirmation
		//   - args: Additional arguments
		// Returns the queue and any error encountered.
		QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)

		// QueueBind binds a queue to an exchange.
		// Parameters:
		//   - name: The name of the queue
		//   - key: The routing key to use
		//   - exchange: The name of the exchange
		//   - noWait: Don't wait for a server confirmation
		//   - args: Additional arguments
		QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error

		// Consume starts delivering messages from a queue.
		// Parameters:
		//   - queue: The name of the queue
		//   - consumer: The consumer tag (empty string to have the server generate one)
		//   - autoAck: Acknowledge messages automatically when delivered
		//   - exclusive: Request exclusive consumer access
		//   - noLocal: Don't deliver messages published on this connection
		//   - noWait: Don't wait for a server confirmation
		//   - args: Additional arguments
		// Returns a channel of delivered messages and any error encountered.
		Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)

		// Publish publishes a message to an exchange.
		// Parameters:
		//   - exchange: The name of the exchange
		//   - key: The routing key to use
		//   - mandatory: Return message if it can't be routed to a queue
		//   - immediate: Return message if it can't be delivered to a consumer immediately
		//   - msg: The message to publish
		Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	}
)

// dial is a variable that holds the function to establish a connection to RabbitMQ.
// It allows for mocking in tests.
var dial = func(cfg *configs.RabbitMQConfigs) (RMQConnection, error) {
	return amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%s", cfg.Schema, cfg.User, cfg.Password, cfg.VHost, cfg.Port))
}

// NewConnection creates a new RabbitMQ connection and channel.
// It establishes a connection to the RabbitMQ server using the provided configuration,
// then creates a channel on that connection.
// Returns the connection, channel, and any error encountered.
func NewConnection(cfgs *configs.Configs) (RMQConnection, AMQPChannel, error) {
	logger := cfgs.Logger

	logger.Debug(LogMessage("connecting to rabbitmq..."))
	conn, err := dial(cfgs.RabbitMQConfigs)
	if err != nil {
		logger.Error(LogMessage("failure to connect to the broker"), zap.Error(err))
		return nil, nil, rabbitMQDialError(err)
	}
	logger.Debug(LogMessage("connected to rabbitmq"))

	logger.Debug(LogMessage("creating amqp channel..."))
	ch, err := conn.Channel()
	if err != nil {
		logger.Error(LogMessage("failure to establish the channel"), zap.Error(err))
		return nil, nil, getChannelError(err)
	}
	logger.Debug(LogMessage("created amqp channel"))

	return conn, ch, nil
}
