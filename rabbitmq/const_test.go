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
	s.Equal("[gokit::rabbitmq] msg", Message("msg"))
}

func (s *RabbitMQConstSuiteTest) TestLogWithType() {
	msg, field := LogMsgWithType("msg", "type", "msgId")

	s.Equal("[gokit::rabbitmq] msgtype", msg)
	s.Equal("messageId", field.Key)
	s.Equal("msgId", field.String)
}

func (s *RabbitMQConstSuiteTest) TestLogMsgWithMessageId() {
	msg, field := LogMsgWithMessageId("msg", "msgId")

	s.Equal("[gokit::rabbitmq] msg", msg)
	s.Equal("messageId", field.Key)
	s.Equal("msgId", field.String)
}
