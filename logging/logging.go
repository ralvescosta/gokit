// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package logging provides structured logging capabilities powered by Zap
// with various configuration options for different environments.
package logging

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger defines the interface for logging operations within the application.
	// It provides methods for different log levels and the ability to add context fields.
	Logger interface {
		// With adds structured context to the logger.
		With(fields ...zapcore.Field) *zap.Logger

		// Debug logs a message at Debug level with optional fields.
		Debug(msg string, fields ...zap.Field)

		// Info logs a message at Info level with optional fields.
		Info(msg string, fields ...zap.Field)

		// Warn logs a message at Warn level with optional fields.
		Warn(msg string, fields ...zap.Field)

		// Error logs a message at Error level with optional fields.
		Error(msg string, fields ...zap.Field)

		// Fatal logs a message at Fatal level with optional fields,
		// then calls os.Exit(1).
		Fatal(msg string, fields ...zap.Field)
	}
)

var (
	// openFile is a variable that holds the os.OpenFile function,
	// allowing it to be replaced in tests.
	openFile = os.OpenFile
)

// NewDefaultLogger creates a new logger that outputs to stdout.
// It configures the logger based on the environment:
// - Production/Staging: Uses JSON encoder
// - Development: Uses colored console output
//
// The log level is determined by the configuration provided.
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

// NewFileLogger creates a logger that outputs to both a file and stdout.
// The file path is specified in the configuration.
// In production/staging environments, it only outputs to the file in JSON format.
// In development environments, it outputs to both stdout (colored) and the file (JSON).
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

// mapZapLogLevel converts the application config log level to the corresponding
// Zap log level. It defaults to InfoLevel if the level is not recognized.
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
