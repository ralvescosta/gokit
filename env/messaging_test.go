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
