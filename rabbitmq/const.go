package rabbitmq

import (
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap/zapcore"
)

const (
	DIRECT_EXCHANGE  ExchangeKind = "direct"
	FANOUT_EXCHANGE  ExchangeKind = "fanout"
	TOPIC_EXCHANGE   ExchangeKind = "topic"
	HEADERS_EXCHANGE ExchangeKind = "headers"
	DELAY_EXCHANGE   ExchangeKind = "x-delayed-message"

	DeclareErrorMessage = "[RabbitMQ::Connect] failure to declare %s: %s"
	BindErrorMessage    = "[RabbitMQ::Connect] failure to bind %s: %s"

	JsonContentType = "application/json"

	AMQPHeaderNumberOfRetry = "x-count"
	AMQPHeaderDelay         = "x-delay"
)

// var (
// 	ErrorConnection               = errors.New("messaging failure to connect to rabbitmq")
// 	ErrorChannel                  = errors.New("messaging error to stablish amqp channel")
// 	ErrorRegisterDispatcher       = errors.New("messaging unformatted dispatcher params")
// 	ErrorRetryable                = errors.New("messaging failure to process send to retry latter")
// 	ErrorReceivedMessageValidator = errors.New("messaging unformatted received message")
// 	ErrorQueueDeclaration         = errors.New("to use dql feature the bind exchanges must be declared first")
// )

func Message(msg string) string {
	return "[gokit::rabbitmq] " + msg
}

func MessageType(msg, typ, msgID string) (string, zapcore.Field) {
	return Message(msg) + typ, logging.MessageIdField(msgID)
}

func MessageId(msg, msgID string) (string, zapcore.Field) {
	return Message(msg), logging.MessageIdField(msgID)
}
