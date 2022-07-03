package rabbitmq

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RabbitMQMessagingSuiteTest struct {
	suite.Suite
}

func TestRabbitMQMessagingSuiteTest(t *testing.T) {
	suite.Run(t, new(RabbitMQMessagingSuiteTest))
}
