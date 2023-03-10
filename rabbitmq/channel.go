package rabbitmq

import "github.com/streadway/amqp"

type (
	AMQPChannel interface {
		ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args map[string]interface{}) error
		ExchangeBind(destination, key, source string, noWait bool, args map[string]interface{}) error
		QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args map[string]interface{}) (amqp.Queue, error)
		QueueBind(name, key, exchange string, noWait bool, args map[string]interface{}) error
		Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args map[string]interface{}) (<-chan amqp.Delivery, error)
		Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	}
)
