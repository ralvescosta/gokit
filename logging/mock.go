package logging

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) With(_ ...zap.Field) *zap.Logger {
	return nil
}
func (m *MockLogger) Debug(_ string, _ ...zap.Field) {
}
func (m *MockLogger) Info(_ string, _ ...zap.Field) {
}
func (m *MockLogger) Warn(_ string, _ ...zap.Field) {
}
func (m *MockLogger) Error(_ string, _ ...zap.Field) {
}
func (m *MockLogger) Fatal(_ string, _ ...zap.Field) {
}

func NewMockLogger() *MockLogger {
	return new(MockLogger)
}
