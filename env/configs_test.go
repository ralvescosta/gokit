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

func (s *EnvTestSuite) SetupTest() {
	dotEnvConfig = func(path string) error {
		return nil
	}
}

func (s *EnvTestSuite) TestNew() {
	os.Setenv("GO_ENV", "dev")

	builder := New()

	s.IsType(&ConfigBuilderImpl{}, builder)
}

func (s *EnvTestSuite) TestAppName() {
	os.Setenv(APP_NAME_ENV_KEY, "app")

	builder := New()

	s.Equal(builder.appName(), "app")
}

func (s *EnvTestSuite) TestLogPath() {
	os.Setenv(LOG_PATH_ENV_KEY, ".")
	builder := New()
	path, _ := os.Getwd()

	s.Contains(builder.logPath(DEFAULT_APP_NAME), path)
}
