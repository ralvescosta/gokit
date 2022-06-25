package mock

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
func (m *MockLogger) Warn(msg string, fields ...zap.Field) {
	// m.Called(msg, fields)
}
func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	// m.Called(msg, fields)
}

func NewMockLogger() *MockLogger {
	return new(MockLogger)
}
