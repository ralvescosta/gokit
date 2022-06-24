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

func (s *MessagingTestSuite) TestGetEngines() {
	os.Setenv(MESSAGING_ENGINES_ENV_KEY, RABBITMQ_ENGINE)
	c := &Configs{}

	c.getEngines()

	s.True(c.MESSAGING_ENGINES[RABBITMQ_ENGINE])
}

func (s *MessagingTestSuite) TestGetEnginesInvalidEngine() {
	os.Setenv(MESSAGING_ENGINES_ENV_KEY, "invalid")
	c := &Configs{}

	c.getEngines()

	s.Error(c.Err)
}

func (s *MessagingTestSuite) TestGetEnginesErr() {
	os.Setenv(MESSAGING_ENGINES_ENV_KEY, "")
	c := &Configs{}

	c.getEngines()

	s.Error(c.Err)
}

func (s *MessagingTestSuite) TestGetRabbitMQConfigs() {
	c := &Configs{
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
	c := &Configs{}
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
	c := &Configs{}

	c.getKafkaConfigs()

	s.Equal(c.KAFKA_HOST, "")
}
