package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HTTPServerTestSuite struct {
	suite.Suite
}

func TestHTTPServerTestSuite(t *testing.T) {
	suite.Run(t, new(HTTPServerTestSuite))
}

func (s *HTTPServerTestSuite) SetupTest() {
	dotEnvConfig = func(path string) error {
		return nil
	}
}

func (s *HTTPServerTestSuite) TestHTTPServerConfigs() {
	os.Setenv(GO_ENV_KEY, "dev")
	os.Setenv(HTTP_PORT_ENV_KEY, "8000")
	os.Setenv(HTTP_HOST_ENV_KEY, "localhost")

	cfg, err := New().HTTPServer().Build()

	s.NoError(err)
	s.NotNil(cfg.HTTPConfigs)

	cfg, err = New().Build()

	s.Nil(err)
	s.Nil(cfg.OtelConfigs)
}

func (s *HTTPServerTestSuite) TestHTTPServerConfigsErr() {
	os.Setenv(GO_ENV_KEY, "dev")

	os.Setenv(HTTP_PORT_ENV_KEY, "")
	_, err := New().HTTPServer().Build()

	s.Error(err)

	//
	os.Setenv(HTTP_PORT_ENV_KEY, "8000")
	os.Setenv(HTTP_HOST_ENV_KEY, "")
	_, err = New().HTTPServer().Build()

	s.Error(err)
}
