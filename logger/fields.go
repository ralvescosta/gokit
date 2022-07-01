package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	MessageIdFieldKey = "messageId"
	AccountIdFieldKey = "accountId"
	ErrorFieldKey     = "error"
)

func uuidField(key string, value string) zap.Field {
	return zap.Field{
		Key:       key,
		Type:      zapcore.StringType,
		Interface: value,
	}
}

func MessageIdField(msgId string) zap.Field {
	return uuidField(MessageIdFieldKey, msgId)
}

func AccountIdField(accId string) zap.Field {
	return uuidField(AccountIdFieldKey, accId)
}

func ErrorField(err error) zap.Field {
	return zap.Field{
		Key:       ErrorFieldKey,
		Type:      zapcore.ErrorType,
		Interface: err,
	}
}
