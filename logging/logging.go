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

var (
	openFile = os.OpenFile
)

func NewDefaultLogger(cfgs *configs.Configs) (Logger, error) {
	zapLogLevel := mapZapLogLevel(cfgs.AppConfigs)

	if cfgs.AppConfigs.GoEnv == configs.ProductionEnv || cfgs.AppConfigs.GoEnv == configs.StagingEnv {
		logConfig := zap.NewProductionEncoderConfig()
		logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(logConfig)

		cfgs.Logger = zap.New(zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLogLevel)).Named(cfgs.AppConfigs.AppName)

		return cfgs.Logger, nil
	}

	logConfig := zap.NewDevelopmentEncoderConfig()
	logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(logConfig)

	cfgs.Logger = zap.New(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapLogLevel)).Named(cfgs.AppConfigs.AppName)

	return cfgs.Logger, nil
}

func NewFileLogger(cfgs *configs.Configs) (Logger, error) {
	zapLogLevel := mapZapLogLevel(cfgs.AppConfigs)

	file, err := openFile(
		cfgs.AppConfigs.LogPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	if cfgs.AppConfigs.GoEnv == configs.ProductionEnv || cfgs.AppConfigs.GoEnv == configs.StagingEnv {
		logConfig := zap.NewProductionEncoderConfig()
		logConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewJSONEncoder(logConfig)

		cfgs.Logger = zap.New(zapcore.NewCore(encoder, zapcore.AddSync(file), zapLogLevel)).Named(cfgs.AppConfigs.AppName)

		return cfgs.Logger, nil
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

	cfgs.Logger = zap.New(core).Named(cfgs.AppConfigs.AppName)

	return cfgs.Logger, nil
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
