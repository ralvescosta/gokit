package env

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrorTestSuite struct {
	suite.Suite
}

func TestErrorTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}

func (s *ErrorTestSuite) SetupTest() {
}

func (s *ErrorTestSuite) TestConfigsError() {
	err := NewConfigsError("some message")

	s.Equal("configs builder error - some message", err.Error())
}

func (s *ErrorTestSuite) TestErrRequiredConfig() {
	err := NewErrRequiredConfig("some message")

	s.Equal("configs builder error - some message is required", err.Error())
}

func (s *ErrorTestSuite) TestError() {
	msg := NewConfigsError("some message").Error()

	s.Equal("configs builder error - some message", msg)
}
