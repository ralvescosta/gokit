package rabbitmq

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RabbitMQConstSuiteTest struct {
	suite.Suite
}

func TestRabbitMQSuiteTest(t *testing.T) {
	suite.Run(t, new(RabbitMQConstSuiteTest))
}

func (s *RabbitMQConstSuiteTest) TestLogMessage() {
	s.Equal(LogMessage("msg"), "[Pkg::RabbitMQ] msg")
}

func (s *RabbitMQConstSuiteTest) TestLogWithType() {
	msg, field := LogMsgWithType("msg", "type", "msgId")

	s.Equal(msg, "[Pkg::RabbitMQ] msgtype")
	s.Equal(field.Key, "messageId")
	s.Equal(field.String, "msgId")
}

func (s *RabbitMQConstSuiteTest) TestLogMsgWithMessageId() {
	msg, field := LogMsgWithMessageId("msg", "msgId")

	s.Equal(msg, "[Pkg::RabbitMQ] msg")
	s.Equal(field.Key, "messageId")
	s.Equal(field.String, "msgId")
}
