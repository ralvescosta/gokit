package logging

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
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

func NewDefaultLogger(e *configs.AppConfigs) (Logger, error) {
	zapLogLevel := mapZapLogLevel(e)

	if e.GoEnv == configs.PRODUCTION_ENV || e.GoEnv == configs.STAGING_ENV {
		logConfig := zap.NewProductionEncoderConfig()
		logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(logConfig)

		return zap.New(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLogLevel)).Named(e.AppName), nil
	}

	logConfig := zap.NewDevelopmentEncoderConfig()
	logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(logConfig)

	return zap.New(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLogLevel)).Named(e.AppName), nil
}

func NewFileLogger(e *configs.AppConfigs) (Logger, error) {
	zapLogLevel := mapZapLogLevel(e)

	file, err := openFile(
		e.LogPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	if e.GoEnv == configs.PRODUCTION_ENV || e.GoEnv == configs.STAGING_ENV {
		logConfig := zap.NewProductionEncoderConfig()
		logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(logConfig)

		return zap.New(zapcore.NewCore(encoder, zapcore.AddSync(file), zapLogLevel)).Named(e.AppName), nil
	}

	logConfig := zap.NewDevelopmentEncoderConfig()
	logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(logConfig)
	fileEncoder := zapcore.NewJSONEncoder(logConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLogLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(file), zapLogLevel),
	)

	return zap.New(core).Named(e.AppName), nil
}

func mapZapLogLevel(e *configs.AppConfigs) zapcore.Level {
	switch e.LogLevel {
	case configs.DEBUG:
		return zap.DebugLevel
	case configs.INFO:
		return zap.InfoLevel
	case configs.WARN:
		return zap.WarnLevel
	case configs.ERROR:
		return zap.ErrorLevel
	case configs.PANIC:
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}
