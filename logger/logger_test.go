package logger

import (
	"testing"

	"github.com/ralvescostati/pkgs/env"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type LoggerTestSuite struct {
	suite.Suite
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (s *LoggerTestSuite) TestMapZapLogLevel() {
	s.Equal(mapZapLogLevel(&env.Env{LOG_LEVEL: env.DEBUG_L}), zap.DebugLevel)
	s.Equal(mapZapLogLevel(&env.Env{LOG_LEVEL: env.INFO_L}), zap.InfoLevel)
	s.Equal(mapZapLogLevel(&env.Env{LOG_LEVEL: env.WARN_L}), zap.WarnLevel)
	s.Equal(mapZapLogLevel(&env.Env{LOG_LEVEL: env.ERROR_L}), zap.ErrorLevel)
	s.Equal(mapZapLogLevel(&env.Env{LOG_LEVEL: env.PANIC_L}), zap.PanicLevel)
}

func (s *LoggerTestSuite) TestNewDefaultLoggerProd() {
	env := &env.Env{
		GO_ENV:    env.PRODUCTION_ENV,
		LOG_LEVEL: env.DEBUG_L,
	}

	logger, err := NewDefaultLogger(env)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewDefaultLoggerDev() {
	env := &env.Env{
		GO_ENV:    env.DEVELOPMENT_ENV,
		LOG_LEVEL: env.DEBUG_L,
	}

	logger, err := NewDefaultLogger(env)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}
