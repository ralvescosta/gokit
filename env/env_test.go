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

func (s *EnvTestSuite) TestNewEnvironment() {
	s.Equal(NewEnvironment("development"), DEVELOPMENT_ENV)
	s.Equal(NewEnvironment("DEVELOPMENT"), DEVELOPMENT_ENV)
	s.Equal(NewEnvironment("dev"), DEVELOPMENT_ENV)
	s.Equal(NewEnvironment("production"), PRODUCTION_ENV)
	s.Equal(NewEnvironment("PRODUCTION"), PRODUCTION_ENV)
	s.Equal(NewEnvironment("prod"), PRODUCTION_ENV)
	s.Equal(NewEnvironment("staging"), STAGING_ENV)
	s.Equal(NewEnvironment("STAGING"), STAGING_ENV)
	s.Equal(NewEnvironment("stg"), STAGING_ENV)
	s.Equal(NewEnvironment("qa"), QA_ENV)
	s.Equal(NewEnvironment("QA"), QA_ENV)
	s.Equal(NewEnvironment("unknown"), UNKNOWN_ENV)
}

func (s *EnvTestSuite) TestNewLogLevel() {
	s.Equal(NewLogLevel("debug"), DEBUG_L)
	s.Equal(NewLogLevel("DEBUG"), DEBUG_L)
	s.Equal(NewLogLevel("warn"), WARN_L)
	s.Equal(NewLogLevel("WARN"), WARN_L)
	s.Equal(NewLogLevel("error"), ERROR_L)
	s.Equal(NewLogLevel("ERROR"), ERROR_L)
	s.Equal(NewLogLevel("panic"), PANIC_L)
	s.Equal(NewLogLevel("PANIC"), PANIC_L)
	s.Equal(NewLogLevel("info"), INFO_L)
}

func (s *EnvTestSuite) TestLoad() {
	os.Setenv("GO_ENV", "dev")

	dotEnvConfig = func(path string) error {
		return nil
	}

	appEnv := NewAppEnvironment()

	err := appEnv.Load()

	s.NoError(err)
}

func (s *EnvTestSuite) TestLoadErr() {
	os.Setenv("GO_ENV", "")
	appEnv := NewAppEnvironment()

	err := appEnv.Load()

	s.Error(err)
}

func (s *EnvTestSuite) TestNewAppEnvironment() {
	s.NotNil(NewAppEnvironment())
}
