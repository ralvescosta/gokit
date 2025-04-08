package logging

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ralvescosta/gokit/configs"
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
	s.Equal(mapZapLogLevel(&configs.AppConfigs{LogLevel: configs.DEBUG}), zap.DebugLevel)
	s.Equal(mapZapLogLevel(&configs.AppConfigs{LogLevel: configs.INFO}), zap.InfoLevel)
	s.Equal(mapZapLogLevel(&configs.AppConfigs{LogLevel: configs.WARN}), zap.WarnLevel)
	s.Equal(mapZapLogLevel(&configs.AppConfigs{LogLevel: configs.ERROR}), zap.ErrorLevel)
	s.Equal(mapZapLogLevel(&configs.AppConfigs{LogLevel: configs.PANIC}), zap.PanicLevel)
}

func (s *LoggerTestSuite) TestNewDefaultLoggerProd() {
	logConfigs := configs.Configs{
		AppConfigs: &configs.AppConfigs{
			GoEnv:    configs.ProductionEnv,
			LogLevel: configs.DEBUG,
		},
	}

	logger, err := NewDefaultLogger(&logConfigs)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewDefaultLoggerDev() {
	logConfigs := configs.Configs{
		AppConfigs: &configs.AppConfigs{
			GoEnv:    configs.DevelopmentEnv,
			LogLevel: configs.DEBUG,
		},
	}

	logger, err := NewDefaultLogger(&logConfigs)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewFileLoggerProd() {
	logConfigs := configs.Configs{
		AppConfigs: &configs.AppConfigs{
			GoEnv:    configs.ProductionEnv,
			LogLevel: configs.DEBUG,
			LogPath:  "./log/file.log",
		},
	}

	fmt.Println(os.Getwd())

	logger, err := NewFileLogger(&logConfigs)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewFileLoggerDev() {
	logConfigs := configs.Configs{
		AppConfigs: &configs.AppConfigs{
			GoEnv:    configs.DevelopmentEnv,
			LogLevel: configs.DEBUG,
			LogPath:  "./log/file.log",
		},
	}

	logger, err := NewFileLogger(&logConfigs)

	s.NoError(err)
	s.IsType(&zap.Logger{}, logger)
}

func (s *LoggerTestSuite) TestNewFileLoggerErrInOpenFile() {
	openFile = func(_ string, flag int, perm os.FileMode) (*os.File, error) {
		return nil, errors.New("some error")
	}

	logConfigs := configs.Configs{
		AppConfigs: &configs.AppConfigs{
			GoEnv:    configs.DevelopmentEnv,
			LogLevel: configs.DEBUG,
			LogPath:  "./log/file.log",
		},
	}

	_, err := NewFileLogger(&logConfigs)

	s.Error(err)
}
