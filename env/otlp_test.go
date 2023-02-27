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
	os.Setenv(GO_ENV_KEY, "dev")
	os.Setenv(TRACING_ENABLED_ENV_KEY, "true")
	os.Setenv(METRICS_ENABLED_ENV_KEY, "true")
	os.Setenv(OTLP_ENDPOINT_ENV_KEY, "endpoint")

	cfg, err := New().Otel().Build()

	s.NoError(err)
	s.NotNil(cfg.OtelConfigs)

	cfg, err = New().Build()

	s.Nil(err)
	s.Nil(cfg.OtelConfigs)
}

func (s *TracingTestSuite) TestTracingErr() {
	os.Setenv(GO_ENV_KEY, "dev")

	os.Setenv(TRACING_ENABLED_ENV_KEY, "")
	_, err := New().Otel().Build()

	s.Error(err)
}
