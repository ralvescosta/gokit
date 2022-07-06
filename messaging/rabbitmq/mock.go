package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type (
	MockRabbitMQMessaging struct {
		mock.Mock
	}

	MockAMQPConnection struct {
		mock.Mock
	}

	MockAMQPChannel struct {
		mock.Mock
	}
)

func (m *MockRabbitMQMessaging) Declare(opts *Topology) IRabbitMQMessaging {
	args := m.Called(opts)

	res := args.Get(0).(IRabbitMQMessaging)

	return res
}

func (m *MockRabbitMQMessaging) ApplyBinds() IRabbitMQMessaging {
	args := m.Called(nil)

	res := args.Get(0).(IRabbitMQMessaging)

	return res
}

func (m *MockRabbitMQMessaging) Publisher(exchange, routingKey string, msg any, opts *PublishOpts) error {
	args := m.Called(exchange, routingKey, msg, opts)

	return args.Error(0)
}

func (m *MockRabbitMQMessaging) Consume() error {
	args := m.Called(nil)

	return args.Error(0)
}

func (m *MockRabbitMQMessaging) RegisterDispatcher(event string, handler ConsumerHandler, t any) error {
	args := m.Called(event, handler, t)

	return args.Error(0)
}

func (m *MockRabbitMQMessaging) Build() (IRabbitMQMessaging, error) {
	args := m.Called(nil)

	res := args.Get(0).(IRabbitMQMessaging)

	return res, args.Error(1)
}

func (m *MockAMQPConnection) Channel() (*amqp.Channel, error) {
	called := m.Called()

	res := called.Get(0).(*amqp.Channel)

	return res, called.Error(1)
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

func NewMockRabbitMQMessaging() *MockRabbitMQMessaging {
	return new(MockRabbitMQMessaging)
}

func NewMockAMQPConnection() *MockAMQPConnection {
	return new(MockAMQPConnection)
}

func NewMockAMQPChannel() *MockAMQPChannel {
	return new(MockAMQPChannel)
}
