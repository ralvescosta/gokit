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
	os.Setenv(TRACING_ENABLED_ENV_KEY, "true")
	os.Setenv(METRICS_ENABLED_ENV_KEY, "true")
	os.Setenv(OTLP_ENDPOINT_ENV_KEY, "endpoint")

	c := &Config{}

	c.Otel()

	s.NoError(c.Err)
}

func (s *TracingTestSuite) TestTracingErr() {
	c := &Config{}
	os.Setenv(TRACING_ENABLED_ENV_KEY, "")

	c.Otel()

	s.Error(c.Err)

	os.Setenv(TRACING_ENABLED_ENV_KEY, "false")
	os.Setenv(METRICS_ENABLED_ENV_KEY, "false")
	os.Setenv(OTLP_ENDPOINT_ENV_KEY, "")

	c.Otel()

	s.Error(c.Err)
}
