// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Unit tests for the guid package.
package guid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// UUIDSuiteTest defines the test suite for UUID functions
type UUIDSuiteTest struct {
	suite.Suite
}

// TestUUIDSuiteTest runs the UUID test suite
func TestUUIDSuiteTest(t *testing.T) {
	suite.Run(t, new(UUIDSuiteTest))
}

// TestUUIDFromString tests the UUIDFromString function
func (s *UUIDSuiteTest) TestUUIDFromString() {
	s.NotEqual(UUIDFromString(uuid.New().String()), uuid.Nil)
	s.Equal(UUIDFromString(""), uuid.Nil)
}

// TestByteSliceFromStringUUID tests the ByteSliceFromStringUUID function
func (s *UUIDSuiteTest) TestByteSliceFromStringUUID() {
	u := uuid.New()
	byt, _ := u.MarshalBinary()

	s.Equal(ByteSliceFromStringUUID(u.String()), byt)
	s.Nil(ByteSliceFromStringUUID(""))
}

// TestUUIDFromByteSlice tests the UUIDFromByteSlice function
func (s *UUIDSuiteTest) TestUUIDFromByteSlice() {
	byt, _ := uuid.New().MarshalBinary()

	s.NotEqual(UUIDFromByteSlice(byt), uuid.Nil)
	s.Equal(UUIDFromByteSlice([]byte{}), uuid.Nil)
}

// TestStringFromByteSliceUUID tests the StringFromByteSliceUUID function
func (s *UUIDSuiteTest) TestStringFromByteSliceUUID() {
	u := uuid.New()
	byt, _ := u.MarshalBinary()

	s.Equal(StringFromByteSliceUUID(byt), u.String())
	s.Equal(StringFromByteSliceUUID([]byte{}), uuid.Nil.String())
}

// TestByteSliceFromUUID tests the ByteSliceFromUUID function
func (s *UUIDSuiteTest) TestByteSliceFromUUID() {
	u := uuid.New()
	byt, _ := u.MarshalBinary()

	s.Equal(ByteSliceFromUUID(u), byt)
}
