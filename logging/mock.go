package logging

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) With(fields ...zap.Field) *zap.Logger {
	return nil
}
func (m *MockLogger) Debug(msg string, fields ...zap.Field) {
}
func (m *MockLogger) Info(msg string, fields ...zap.Field) {
}
func (m *MockLogger) Warn(msg string, fields ...zap.Field) {
}
func (m *MockLogger) Error(msg string, fields ...zap.Field) {
}
func (m *MockLogger) Fatal(msg string, fields ...zap.Field) {
}

func NewMockLogger() *MockLogger {
	return new(MockLogger)
}
