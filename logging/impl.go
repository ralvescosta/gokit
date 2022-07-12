package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ralvescosta/toolkit/env"
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
)

var (
	openFile = os.OpenFile
)

func NewDefaultLogger(e *env.Configs) (ILogger, error) {
	zapLogLevel := mapZapLogLevel(e)

	if e.GO_ENV == env.PRODUCTION_ENV || e.GO_ENV == env.STAGING_ENV {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(config)

		return zap.New(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLogLevel)).Named(e.APP_NAME), nil
	}

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	return zap.New(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLogLevel)).Named(e.APP_NAME), nil
}

func NewFileLogger(e *env.Configs) (ILogger, error) {
	zapLogLevel := mapZapLogLevel(e)

	file, err := openFile(
		e.LOG_PATH,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	if e.GO_ENV == env.PRODUCTION_ENV || e.GO_ENV == env.STAGING_ENV {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(config)

		return zap.New(zapcore.NewCore(encoder, zapcore.AddSync(file), zapLogLevel)).Named(e.APP_NAME), nil
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

	return zap.New(core).Named(e.APP_NAME), nil
}

func mapZapLogLevel(e *env.Configs) zapcore.Level {
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
