package logging

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ralvescosta/gokit/env"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type LoggerTestSuite struct {
	suite.Suite
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (s *LoggerTestSuite) SetupTest() {
	openFile = os.OpenFile
}

func (s *LoggerTestSuite) TestMapZapLogLevel() {
	s.Equal(mapZapLogLevel(&env.Config{LOG_LEVEL: env.DEBUG_L}), zap.DebugLevel)
	s.Equal(mapZapLogLevel(&env.Config{LOG_LEVEL: env.INFO_L}), zap.InfoLevel)
	s.Equal(mapZapLogLevel(&env.Config{LOG_LEVEL: env.WARN_L}), zap.WarnLevel)
	s.Equal(mapZapLogLevel(&env.Config{LOG_LEVEL: env.ERROR_L}), zap.ErrorLevel)
	s.Equal(mapZapLogLevel(&env.Config{LOG_LEVEL: env.PANIC_L}), zap.PanicLevel)
}

func (s *LoggerTestSuite) TestNewDefaultLoggerProd() {
	env := &env.Config{
		GO_ENV:    env.PRODUCTION_ENV,
		LOG_LEVEL: env.DEBUG_L,
	}

	logger, err := NewDefaultLogger(env)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewDefaultLoggerDev() {
	env := &env.Config{
		GO_ENV:    env.DEVELOPMENT_ENV,
		LOG_LEVEL: env.DEBUG_L,
	}

	logger, err := NewDefaultLogger(env)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewFileLoggerProd() {
	env := &env.Config{
		GO_ENV:    env.PRODUCTION_ENV,
		LOG_LEVEL: env.DEBUG_L,
		LOG_PATH:  "./log/file.log",
	}

	fmt.Println(os.Getwd())

	logger, err := NewFileLogger(env)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewFileLoggerDev() {
	env := &env.Config{
		GO_ENV:    env.DEVELOPMENT_ENV,
		LOG_LEVEL: env.DEBUG_L,
		LOG_PATH:  "./log/file.log",
	}

	logger, err := NewFileLogger(env)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewFileLoggerErrInOpenFile() {
	openFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, errors.New("some error")
	}

	env := &env.Config{
		GO_ENV:    env.DEVELOPMENT_ENV,
		LOG_LEVEL: env.DEBUG_L,
		LOG_PATH:  "./log/file.log",
	}

	_, err := NewFileLogger(env)

	s.Error(err)
}
