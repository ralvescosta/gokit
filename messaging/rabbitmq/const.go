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
	DELAY_EXCHANGE   ExchangeKind = "delay"

	ConnErrorMessage    = "[RabbitMQ::Connect] failure to connect to the %s: %s"
	DeclareErrorMessage = "[RabbitMQ::Connect] failure to declare %s: %s"
	BindErrorMessage    = "[RabbitMQ::Connect] failure to bind %s: %s"

	DeadLetterSuffix = "-dead-letter"
	JsonContentType  = "application/json"

	AMQPHeaderNumberOfRetry   = "x-count"
	AMQPHeaderTraceID         = "x-trace-id"
	AMQPHeaderRejected        = "x-rejected"
	AMQPHeaderRejectionReason = "x-rejection-reason"
)

func LogMessageId(msgID string) zap.Field {
	return zap.Field{
		Key:    "MessageId",
		Type:   zapcore.StringType,
		String: msgID,
	}
}

func LogMessage(msg string) string {
	return "[RabbitMQ:HandlerExecutor] " + msg
}

func LogMsgWithType(msg, typ, msgID string) (string, zap.Field) {
	return LogMessage(msg) + typ, LogMessageId(msgID)
}

func LogMsgWithMessageId(msg, msgID string) (string, zap.Field) {
	return LogMessage(msg), LogMessageId(msgID)
}
