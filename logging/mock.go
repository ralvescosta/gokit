// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package logging

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockLogger is a mock implementation of the Logger interface
// that can be used in unit tests to verify logging behavior
// without actually producing log output.
type MockLogger struct {
	mock.Mock
}

// With implements the Logger interface's With method for the mock.
// Returns nil instead of a real logger since it's a mock.
func (m *MockLogger) With(_ ...zap.Field) *zap.Logger {
	return nil
}

// Debug implements the Logger interface's Debug method for the mock.
func (m *MockLogger) Debug(_ string, _ ...zap.Field) {
}

// Info implements the Logger interface's Info method for the mock.
func (m *MockLogger) Info(_ string, _ ...zap.Field) {
}

// Warn implements the Logger interface's Warn method for the mock.
func (m *MockLogger) Warn(_ string, _ ...zap.Field) {
}

// Error implements the Logger interface's Error method for the mock.
func (m *MockLogger) Error(_ string, _ ...zap.Field) {
}

// Fatal implements the Logger interface's Fatal method for the mock.
func (m *MockLogger) Fatal(_ string, _ ...zap.Field) {
}

// NewMockLogger creates and returns a new instance of MockLogger
// that can be used in tests.
func NewMockLogger() *MockLogger {
	return new(MockLogger)
}
