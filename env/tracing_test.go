package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TracingTestSuite struct {
	suite.Suite
}

func TestTracingTestSuite(t *testing.T) {
	suite.Run(t, new(TracingTestSuite))
}

func (s *TracingTestSuite) SetupTest() {
	dotEnvConfig = func(path string) error {
		return nil
	}
}

func (s *TracingTestSuite) TestTracing() {
	os.Setenv(IS_TRACING_ENABLED_ENV_KEY, "true")
	os.Setenv(OTLP_ENDPOINT_ENV_KEY, "endpoint")

	c := &Config{}

	c.Tracing()

	s.NoError(c.Err)
}

func (s *TracingTestSuite) TestTracingErr() {
	c := &Config{}
	os.Setenv(IS_TRACING_ENABLED_ENV_KEY, "")

	c.Tracing()
	s.Error(c.Err)

	os.Setenv(IS_TRACING_ENABLED_ENV_KEY, "false")
	os.Setenv(OTLP_ENDPOINT_ENV_KEY, "")

	c.Tracing()
	s.Error(c.Err)
}
