package rabbitmq

import (
	"fmt"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type (
	AMQPChannel interface {
		ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
		ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error
		QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
		QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
		Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
		Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	}
)

var dial = func(cfg *env.RabbitMQConfigs) (AMQPConnection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.VHost, cfg.Port))
}

func NewChannel(cfg *env.RabbitMQConfigs, logger logging.Logger) (AMQPChannel, error) {
	logger.Debug(LogMessage("connecting to rabbitmq..."))
	conn, err := dial(cfg)
	if err != nil {
		logger.Error(LogMessage("failure to connect to the broker"), zap.Error(err))
		return nil, err
	}
	logger.Debug(LogMessage("connected to rabbitmq"))

	logger.Debug(LogMessage("creating amqp channel..."))
	ch, err := conn.Channel()
	if err != nil {
		logger.Error(LogMessage("failure to establish the channel"), zap.Error(err))
		return nil, err
	}
	logger.Debug(LogMessage("created amqp channel"))

	return ch, nil
}
