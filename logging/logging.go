package logging

import (
	"os"

	"github.com/ralvescosta/gokit/env"
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

var (
	openFile = os.OpenFile
)

func NewDefaultLogger(e *env.AppConfigs) (Logger, error) {
	zapLogLevel := mapZapLogLevel(e)

	if e.GoEnv == env.PRODUCTION_ENV || e.GoEnv == env.STAGING_ENV {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(config)

		return zap.New(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLogLevel)).Named(e.AppName), nil
	}

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	return zap.New(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLogLevel)).Named(e.AppName), nil
}

func NewFileLogger(e *env.AppConfigs) (Logger, error) {
	zapLogLevel := mapZapLogLevel(e)

	file, err := openFile(
		e.LogPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	if e.GoEnv == env.PRODUCTION_ENV || e.GoEnv == env.STAGING_ENV {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(config)

		return zap.New(zapcore.NewCore(encoder, zapcore.AddSync(file), zapLogLevel)).Named(e.AppName), nil
	}

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	fileEncoder := zapcore.NewJSONEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLogLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zapLogLevel),
	)

	return zap.New(core).Named(e.AppName), nil
}

func mapZapLogLevel(e *env.AppConfigs) zapcore.Level {
	switch e.LogLevel {
	case env.DEBUG_L:
		return zap.DebugLevel
	case env.INFO_L:
		return zap.InfoLevel
	case env.WARN_L:
		return zap.WarnLevel
	case env.ERROR_L:
		return zap.ErrorLevel
	case env.PANIC_L:
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}
