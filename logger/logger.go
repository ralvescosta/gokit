package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ralvescostati/pkgs/env"
	"github.com/stretchr/testify/mock"
)

type (
	ILogger interface {
		Debug(msg string, fields ...zap.Field)
		Info(msg string, fields ...zap.Field)
		Warn(msg string, fields ...zap.Field)
		Error(msg string, fields ...zap.Field)
	}

	LogField struct {
		Key   string
		Value interface{}
	}

	MockLogger struct {
		mock.Mock
	}
)

func NewDefaultLogger(e *env.Env) (ILogger, error) {
	zapLogLevel := mapZapLogLevel(e)

	if e.GO_ENV == env.PRODUCTION_ENV || e.GO_ENV == env.STAGING_ENV {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(config)

		return zap.New(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLogLevel)), nil
	}

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	return zap.New(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLogLevel)), nil
}

func NewFileLogger(e *env.Env) (ILogger, error) {
	return nil, nil
}

func mapZapLogLevel(e *env.Env) zapcore.Level {
	switch e.LOG_LEVEL {
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
