package rabbitmq

import (
	"errors"
	"testing"

	// "github.com/ralvescostati/pkgs/messaging/rabbitmq/mock"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logging"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/suite"
)

type RabbitMQMessagingSuiteTest struct {
	suite.Suite

	amqpConn    *MockAMQPConnection
	amqpConnErr error
	amqpChannel *MockAMQPChannel
}

func TestRabbitMQMessagingSuiteTest(t *testing.T) {
	suite.Run(t, new(RabbitMQMessagingSuiteTest))
}

func (s *RabbitMQMessagingSuiteTest) SetupTest() {
	s.amqpConn = NewMockAMQPConnection()
	s.amqpConnErr = nil
	s.amqpChannel = NewMockAMQPChannel()

	dial = func(cfg *env.Configs) (AMQPConnection, error) {
		return s.amqpConn, s.amqpConnErr
	}
}

func (s *RabbitMQMessagingSuiteTest) TestNew() {
	s.amqpConn.
		On("Channel").
		Return(&amqp.Channel{}, nil)

	msg := New(&env.Configs{}, logging.NewMockLogger())
	conn, err := msg.Build()

	s.NotNil(conn)
	s.NoError(err)
}

func (s *RabbitMQMessagingSuiteTest) TestNewConnErr() {
	s.amqpConnErr = errors.New("some err")

	msg := New(&env.Configs{}, logging.NewMockLogger())
	conn, err := msg.Build()

	s.Nil(conn)
	s.Error(err)
}

func (s *RabbitMQMessagingSuiteTest) TestNewChannelErr() {
	s.amqpConn.
		On("Channel").
		Return(&amqp.Channel{}, errors.New("some error"))

	msg := New(&env.Configs{}, logging.NewMockLogger())
	conn, err := msg.Build()

	s.Nil(conn)
	s.Error(err)
}
