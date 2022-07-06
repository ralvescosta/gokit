package uuid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UUIDSuiteTest struct {
	suite.Suite
}

func TestUUIDSuiteTest(t *testing.T) {
	suite.Run(t, new(UUIDSuiteTest))
}

func (s *UUIDSuiteTest) TestUUIDFromString() {
	s.NotEqual(UUIDFromString(uuid.NewString()), uuid.Nil)
	s.Equal(UUIDFromString(""), uuid.Nil)
}

func (s *UUIDSuiteTest) TestByteSliceFromStringUUID() {
	u := uuid.New()
	byt, _ := u.MarshalBinary()

	s.Equal(ByteSliceFromStringUUID(u.String()), byt)
	s.Nil(ByteSliceFromStringUUID(""))
}

func (s *UUIDSuiteTest) TestUUIDFromByteSlice() {
	byt, _ := uuid.New().MarshalBinary()

	s.NotEqual(UUIDFromByteSlice(byt), uuid.Nil)
	s.Equal(UUIDFromByteSlice([]byte{}), uuid.Nil)
}

func (s *UUIDSuiteTest) TestStringFromByteSliceUUID() {
	u := uuid.New()
	byt, _ := u.MarshalBinary()

	s.Equal(StringFromByteSliceUUID(byt), u.String())
	s.Equal(StringFromByteSliceUUID([]byte{}), uuid.Nil.String())
}

func (s *UUIDSuiteTest) TestByteSliceFromUUID() {
	u := uuid.New()
	byt, _ := u.MarshalBinary()

	s.Equal(ByteSliceFromUUID(u), byt)
}
