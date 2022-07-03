package mock

import (
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

type MockAMQPConnection struct {
	mock.Mock
}

func (m *MockAMQPConnection) Channel() (*amqp.Channel, error) {
	called := m.Called(nil)

	res := called.Get(0).(*amqp.Channel)

	return res, called.Error(1)
}
