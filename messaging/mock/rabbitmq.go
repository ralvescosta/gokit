package mock

import (
	"github.com/ralvescostati/pkgs/messaging/rabbitmq"

	"github.com/stretchr/testify/mock"
)

type MockRabbitMQMessaging struct {
	mock.Mock
}

func (m *MockRabbitMQMessaging) Declare(opts *rabbitmq.Topology) rabbitmq.IRabbitMQMessaging {
	args := m.Called(opts)

	res := args.Get(0).(rabbitmq.IRabbitMQMessaging)

	return res
}

func (m *MockRabbitMQMessaging) ApplyBinds() rabbitmq.IRabbitMQMessaging {
	args := m.Called(nil)

	res := args.Get(0).(rabbitmq.IRabbitMQMessaging)

	return res
}

func (m *MockRabbitMQMessaging) Publisher(exchange, routingKey string, msg any, opts *rabbitmq.PublishOpts) error {
	args := m.Called(exchange, routingKey, msg, opts)

	return args.Error(0)
}

func (m *MockRabbitMQMessaging) Consume() error {
	args := m.Called(nil)

	return args.Error(0)
}

func (m *MockRabbitMQMessaging) RegisterDispatcher(event string, handler rabbitmq.ConsumerHandler, t any) error {
	args := m.Called(event, handler, t)

	return args.Error(0)
}

func (m *MockRabbitMQMessaging) Build() (rabbitmq.IRabbitMQMessaging, error) {
	args := m.Called(nil)

	res := args.Get(0).(rabbitmq.IRabbitMQMessaging)

	return res, args.Error(1)
}

func NewMockRabbitMQMessaging() *MockRabbitMQMessaging {
	return new(MockRabbitMQMessaging)
}
