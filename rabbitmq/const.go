package rabbitmq

import (
	"go.uber.org/zap"
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

func Message(msg string) string {
	return "[gokit::rabbitmq] " + msg
}

func MessageType(msg, typ, msgID string) (string, zapcore.Field) {
	return Message(msg) + typ, zap.String("messageId", msgID)
}

func MessageId(msg, msgID string) (string, zapcore.Field) {
	return Message(msg), zap.String("messageId", msgID)
}
