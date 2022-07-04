package env

import "github.com/stretchr/testify/mock"

type MockEnv struct {
	mock.Mock
}

func (m *MockEnv) Load() error {
	called := m.Called()
	return called.Error(0)
}

func NewMockEnv() *MockEnv {
	return new(MockEnv)
}
