package logger

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	MessageIdFieldKey = "messageId"
	AccountIdFieldKey = "accountId"
	ErrorFieldKey     = "error"
)

type UUID interface {
	string | []byte | uuid.UUID
}

func uuidField[T UUID](key string, value T) zap.Field {
	return zap.Field{
		Key:       key,
		Type:      zapcore.UnknownType,
		Interface: value,
	}
}

func MessageIdField[T UUID](msgId T) zap.Field {
	return uuidField(MessageIdFieldKey, msgId)
}

func AccountIdField[T UUID](accId T) zap.Field {
	return uuidField(AccountIdFieldKey, accId)
}

func ErrorField(err error) zap.Field {
	return zap.Field{
		Key:       ErrorFieldKey,
		Type:      zapcore.ErrorType,
		Interface: err,
	}
}
