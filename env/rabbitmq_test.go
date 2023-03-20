package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MessagingTestSuite struct {
	suite.Suite
}

func TestMessagingTestSuite(t *testing.T) {
	suite.Run(t, new(MessagingTestSuite))
}

func (s *MessagingTestSuite) SetupTest() {
	dotEnvConfig = func(path string) error {
		return nil
	}
}

func (s *MessagingTestSuite) TestRabbitMQ() {
	os.Setenv(GO_ENV_KEY, "dev")
	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	os.Setenv(RABBIT_PORT_ENV_KEY, "port")
	os.Setenv(RABBIT_USER_ENV_KEY, "user")
	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "password")
	os.Setenv(RABBIT_VHOST_ENV_KEY, "/")

	cfg, err := New().RabbitMQ().Build()

	s.NoError(err)
	s.NotNil(cfg.RabbitMQConfigs)

	cfg, err = New().Build()
	s.Nil(cfg.RabbitMQConfigs)
	s.Nil(err)
}

func (s *MessagingTestSuite) TestRabbitMQConfigsErr() {
	os.Setenv(GO_ENV_KEY, "dev")

	//
	os.Setenv(RABBIT_HOST_ENV_KEY, "")

	_, err := New().RabbitMQ().Build()

	s.Error(err)

	//
	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	os.Setenv(RABBIT_PORT_ENV_KEY, "")

	_, err = New().RabbitMQ().Build()

	s.Error(err)

	//
	os.Setenv(RABBIT_PORT_ENV_KEY, "port")
	os.Setenv(RABBIT_USER_ENV_KEY, "")

	_, err = New().RabbitMQ().Build()

	s.Error(err)

	//
	os.Setenv(RABBIT_USER_ENV_KEY, "user")
	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "")

	_, err = New().RabbitMQ().Build()

	s.Error(err)
}
