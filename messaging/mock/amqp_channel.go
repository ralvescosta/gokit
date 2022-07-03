package mock

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type MockAMQPChannel struct {
	mock.Mock
}

func (m *MockAMQPChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	called := m.Called(name, kind, durable, autoDelete, internal, noWait, args)

	return called.Error(0)
}

func (m *MockAMQPChannel) ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error {
	called := m.Called(destination, key, source, noWait, args)

	return called.Error(0)
}

func (m *MockAMQPChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	called := m.Called(name, durable, autoDelete, exclusive, noWait, args)

	res := called.Get(0).(amqp.Queue)

	return res, called.Error(1)
}

func (m *MockAMQPChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	called := m.Called(name, key, exchange, noWait, args)

	return called.Error(0)
}

func (m *MockAMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	called := m.Called(queue, consumer, autoAck, exclusive, noLocal, noWait, args)

	res := called.Get(0).(<-chan amqp.Delivery)

	return res, called.Error(1)
}

func (m *MockAMQPChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	called := m.Called(exchange, key, mandatory, immediate, msg)

	return called.Error(0)
}
