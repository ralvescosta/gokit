package rabbitmq

import (
	"errors"

	"github.com/ralvescostati/pkgs/logging"
	"go.uber.org/zap/zapcore"
)

const (
	DIRECT_EXCHANGE  ExchangeKind = "direct"
	FANOUT_EXCHANGE  ExchangeKind = "fanout"
	TOPIC_EXCHANGE   ExchangeKind = "topic"
	HEADERS_EXCHANGE ExchangeKind = "headers"
	DELAY_EXCHANGE   ExchangeKind = "x-delayed-message"

	DLQ_FALLBACK   FallbackType = "dlq"
	RETRY_FALLBACK FallbackType = "delayed"

	DeclareErrorMessage = "[RabbitMQ::Connect] failure to declare %s: %s"
	BindErrorMessage    = "[RabbitMQ::Connect] failure to bind %s: %s"

	JsonContentType = "application/json"

	AMQPHeaderNumberOfRetry = "x-count"
	AMQPHeaderTraceID       = "x-trace-id"
	AMQPHeaderDelay         = "x-delay"
)

var (
	ErrorConnection               = errors.New("messaging failure to connect to rabbitmq")
	ErrorChannel                  = errors.New("messaging error to stablish amqp channel")
	ErrorRegisterDispatcher       = errors.New("messaging unformatted dispatcher params")
	ErrorRetryable                = errors.New("messaging failure to process send to retry latter")
	ErrorReceivedMessageValidator = errors.New("messaging unformatted received message")
)

func LogMessage(msg string) string {
	return "[Pkg::RabbitMQ] " + msg
}

func LogMsgWithType(msg, typ string, msgID string) (string, zapcore.Field) {
	return LogMessage(msg) + typ, logging.MessageIdField(msgID)
}

func LogMsgWithMessageId(msg string, msgID string) (string, zapcore.Field) {
	return LogMessage(msg), logging.MessageIdField(msgID)
}
