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

func (s *MessagingTestSuite) TestMessaging() {
	c := &Config{}
	os.Setenv(MESSAGING_ENGINES_ENV_KEY, RABBITMQ_ENGINE+","+KAFKA_ENGINE)
	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	os.Setenv(RABBIT_PORT_ENV_KEY, "port")
	os.Setenv(RABBIT_USER_ENV_KEY, "user")
	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "password")
	os.Setenv(RABBIT_VHOST_ENV_KEY, "/")

	c.Messaging()

	s.NoError(c.Err)
}

func (s *MessagingTestSuite) TestMessagingErr() {
	c := &Config{}
	c.Err = errors.New("some error")

	s.Nil(c.MESSAGING_ENGINES)

	c.Err = nil
	os.Setenv(MESSAGING_ENGINES_ENV_KEY, "")
	c.Messaging()

	s.Error(c.Err)

	c.Err = nil
	os.Setenv(MESSAGING_ENGINES_ENV_KEY, RABBITMQ_ENGINE+","+KAFKA_ENGINE)
	os.Setenv(RABBIT_HOST_ENV_KEY, "")
	c.Messaging()

	s.Error(c.Err)

	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	os.Setenv(RABBIT_PORT_ENV_KEY, "port")
	os.Setenv(RABBIT_USER_ENV_KEY, "user")
	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "password")
	os.Setenv(RABBIT_VHOST_ENV_KEY, "/")
	c.Messaging()

	// s.Error(c.Err)
}

func (s *MessagingTestSuite) TestGetEngines() {
	os.Setenv(MESSAGING_ENGINES_ENV_KEY, RABBITMQ_ENGINE+","+KAFKA_ENGINE)
	c := &Config{}

	c.getEngines()

	s.True(c.MESSAGING_ENGINES[RABBITMQ_ENGINE])
}

func (s *MessagingTestSuite) TestGetEnginesInvalidEngine() {
	os.Setenv(MESSAGING_ENGINES_ENV_KEY, "invalid")
	c := &Config{}

	c.getEngines()

	s.Error(c.Err)
}

func (s *MessagingTestSuite) TestGetEnginesErr() {
	os.Setenv(MESSAGING_ENGINES_ENV_KEY, "")
	c := &Config{}

	c.getEngines()

	s.Error(c.Err)
}

func (s *MessagingTestSuite) TestGetRabbitMQConfigs() {
	c := &Config{
		MESSAGING_ENGINES: map[string]bool{RABBITMQ_ENGINE: true},
	}
	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	os.Setenv(RABBIT_PORT_ENV_KEY, "port")
	os.Setenv(RABBIT_USER_ENV_KEY, "user")
	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "password")
	os.Setenv(RABBIT_VHOST_ENV_KEY, "/")

	c.getRabbitMQConfigs()

	s.Equal(c.RABBIT_HOST, "host")
	s.Equal(c.RABBIT_PORT, "port")
	s.Equal(c.RABBIT_USER, "user")
	s.Equal(c.RABBIT_PASSWORD, "password")
	s.Equal(c.RABBIT_VHOST, "/")
}

func (s *MessagingTestSuite) TestGetRabbitMQConfigsErr() {
	c := &Config{}
	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	c.getRabbitMQConfigs()

	s.Equal(c.RABBIT_HOST, "")

	c.MESSAGING_ENGINES = map[string]bool{RABBITMQ_ENGINE: true}
	os.Setenv(RABBIT_HOST_ENV_KEY, "")
	c.getRabbitMQConfigs()

	s.Error(c.Err)

	os.Setenv(RABBIT_HOST_ENV_KEY, "host")
	os.Setenv(RABBIT_PORT_ENV_KEY, "")
	c.getRabbitMQConfigs()

	s.Error(c.Err)

	os.Setenv(RABBIT_PORT_ENV_KEY, "port")
	os.Setenv(RABBIT_USER_ENV_KEY, "")
	c.getRabbitMQConfigs()

	s.Error(c.Err)

	os.Setenv(RABBIT_USER_ENV_KEY, "user")
	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "")
	c.getRabbitMQConfigs()

	s.Error(c.Err)

	os.Setenv(RABBIT_PASSWORD_ENV_KEY, "password")
	os.Setenv(RABBIT_VHOST_ENV_KEY, "")
	c.getRabbitMQConfigs()

	s.Error(c.Err)
}

func (s *MessagingTestSuite) TestTetKafkaConfigs() {
	c := &Config{}

	c.getKafkaConfigs()

	s.Equal(c.KAFKA_HOST, "")
}
