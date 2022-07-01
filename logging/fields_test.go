package logging

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type FieldsTestSuite struct {
	suite.Suite
}

func TestFieldsTestSuite(t *testing.T) {
	suite.Run(t, new(FieldsTestSuite))
}

func (s *FieldsTestSuite) TestMessageIdField() {
	f := MessageIdField("uuid")

	s.Equal(f.Key, MessageIdFieldKey)
	s.IsType(zap.Field{}, f)
}

func (s *FieldsTestSuite) TestAccountIdField() {
	f := AccountIdField("uuid")

	s.Equal(f.Key, AccountIdFieldKey)
	s.IsType(zap.Field{}, f)
}

func (s *FieldsTestSuite) TestErrorField() {
	f := ErrorField(errors.New("some error"))

	s.Equal(f.Key, ErrorFieldKey)
	s.IsType(zap.Field{}, f)
}
