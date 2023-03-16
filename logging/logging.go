package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger interface {
		With(fields ...zapcore.Field) *zap.Logger
		Debug(msg string, fields ...zap.Field)
		Info(msg string, fields ...zap.Field)
		Warn(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
		Fatal(msg string, fields ...zap.Field)
	}
)

const (
	MessageIdFieldKey = "messageId"
	AccountIdFieldKey = "accountId"
	ErrorFieldKey     = "error"
)
