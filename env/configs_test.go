package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EnvTestSuite struct {
	suite.Suite
}

func TestEnvTestSuite(t *testing.T) {
	suite.Run(t, new(EnvTestSuite))
}

func (s *EnvTestSuite) TestNew() {
	os.Setenv("GO_ENV", "dev")

	dotEnvConfig = func(path string) error {
		return nil
	}

	builder := New()
	cfg, _ := builder.Build()

	s.Equal(cfg.GO_ENV, DEVELOPMENT_ENV)
	s.NoError(cfg.Err)
}

func (s *EnvTestSuite) TestNewErr() {
	os.Setenv("GO_ENV", "")

	dotEnvConfig = func(path string) error {
		return nil
	}

	builder := New()
	cfg, err := builder.Build()

	s.Equal(cfg.Err.Error(), err.Error())
}

func (s *EnvTestSuite) TestNewAppName() {
	os.Setenv(APP_NAME_ENV_KEY, "")
	s.Equal(NewAppName(), DEFAULT_APP_NAME)

	os.Setenv(APP_NAME_ENV_KEY, "test")
	s.Equal(NewAppName(), "test")
}

func (s *EnvTestSuite) TestNewLogPath() {
	os.Setenv(LOG_PATH_ENV_KEY, "")
	s.Contains(NewLogPath(DEFAULT_APP_NAME), DEFAULT_LOG_PATH)

	path, _ := os.Getwd()

	os.Setenv(LOG_PATH_ENV_KEY, ".")
	s.Contains(NewLogPath(DEFAULT_APP_NAME), path)

	os.Setenv(LOG_PATH_ENV_KEY, "some")
	s.Contains(NewLogPath(DEFAULT_APP_NAME), "/some")
}
