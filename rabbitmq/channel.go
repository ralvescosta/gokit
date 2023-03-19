package rabbitmq

import "github.com/streadway/amqp"

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
