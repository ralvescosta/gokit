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
	s.Equal(mapZapLogLevel(&env.AppConfigs{LogLevel: env.DEBUG_L}), zap.DebugLevel)
	s.Equal(mapZapLogLevel(&env.AppConfigs{LogLevel: env.INFO_L}), zap.InfoLevel)
	s.Equal(mapZapLogLevel(&env.AppConfigs{LogLevel: env.WARN_L}), zap.WarnLevel)
	s.Equal(mapZapLogLevel(&env.AppConfigs{LogLevel: env.ERROR_L}), zap.ErrorLevel)
	s.Equal(mapZapLogLevel(&env.AppConfigs{LogLevel: env.PANIC_L}), zap.PanicLevel)
}

func (s *LoggerTestSuite) TestNewDefaultLoggerProd() {
	env := &env.AppConfigs{
		GoEnv:    env.PRODUCTION_ENV,
		LogLevel: env.DEBUG_L,
	}

	logger, err := NewDefaultLogger(env)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewDefaultLoggerDev() {
	env := &env.AppConfigs{
		GoEnv:    env.DEVELOPMENT_ENV,
		LogLevel: env.DEBUG_L,
	}

	logger, err := NewDefaultLogger(env)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewFileLoggerProd() {
	env := &env.AppConfigs{
		GoEnv:    env.PRODUCTION_ENV,
		LogLevel: env.DEBUG_L,
		LogPath:  "./log/file.log",
	}

	fmt.Println(os.Getwd())

	logger, err := NewFileLogger(env)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewFileLoggerDev() {
	env := &env.AppConfigs{
		GoEnv:    env.DEVELOPMENT_ENV,
		LogLevel: env.DEBUG_L,
		LogPath:  "./log/file.log",
	}

	logger, err := NewFileLogger(env)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewFileLoggerErrInOpenFile() {
	openFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, errors.New("some error")
	}

	env := &env.AppConfigs{
		GoEnv:    env.DEVELOPMENT_ENV,
		LogLevel: env.DEBUG_L,
		LogPath:  "./log/file.log",
	}

	_, err := NewFileLogger(env)

	s.Error(err)
}
