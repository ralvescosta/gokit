package env

import (
	"errors"
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
	c := &Config{}

	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	os.Setenv(RABBIT_PORT_ENV_KEY, "port")
	os.Setenv(RABBIT_USER_ENV_KEY, "user")
	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "password")
	os.Setenv(RABBIT_VHOST_ENV_KEY, "/")

	c.RabbitMQ()

	s.NoError(c.Err)
}

func (s *MessagingTestSuite) TestRabbitMQErr() {
	c := &Config{}
	c.Err = errors.New("some error")

	s.Nil(c.MESSAGING_ENGINES)

	s.Error(c.Err)

	c.Err = nil
	os.Setenv(RABBIT_HOST_ENV_KEY, "")
	c.RabbitMQ()

	s.Error(c.Err)

	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	os.Setenv(RABBIT_PORT_ENV_KEY, "port")
	os.Setenv(RABBIT_USER_ENV_KEY, "user")
	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "password")
	os.Setenv(RABBIT_VHOST_ENV_KEY, "/")
	c.RabbitMQ()

	// s.Error(c.Err)
}

func (s *MessagingTestSuite) TestRabbitMQConfigsErr() {
	c := &Config{}
	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	c.RabbitMQ()

	s.Equal(c.RABBIT_HOST, "host")

	os.Setenv(RABBIT_HOST_ENV_KEY, "")
	c.RabbitMQ()

	s.Error(c.Err)

	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	os.Setenv(RABBIT_PORT_ENV_KEY, "")
	c.RabbitMQ()

	s.Error(c.Err)

	os.Setenv(RABBIT_PORT_ENV_KEY, "port")
	os.Setenv(RABBIT_USER_ENV_KEY, "")
	c.RabbitMQ()

	s.Error(c.Err)

	os.Setenv(RABBIT_USER_ENV_KEY, "user")
	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "")
	c.RabbitMQ()

	s.Error(c.Err)

	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "password")
	os.Setenv(RABBIT_VHOST_ENV_KEY, "")
	c.RabbitMQ()

	s.Error(c.Err)
}
