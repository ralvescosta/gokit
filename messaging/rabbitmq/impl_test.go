package rabbitmq

import (
	"testing"

	// "github.com/ralvescostati/pkgs/messaging/rabbitmq/mock"

	"github.com/stretchr/testify/suite"
)

type RabbitMQMessagingSuiteTest struct {
	suite.Suite

	amqpConn AMQPConnection
}

func TestRabbitMQMessagingSuiteTest(t *testing.T) {
	suite.Run(t, new(RabbitMQMessagingSuiteTest))
}

func (s *RabbitMQMessagingSuiteTest) TestSetup() {
	s.amqpConn = NewMockAMQPConnection()
	// dial = func(cfg *env.Configs) (AMQPConnection, error) {}
}
