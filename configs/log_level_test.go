// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewLogLevel verifies that the NewLogLevel function correctly converts
// string representations of log levels into their corresponding LogLevel enum values.
// It tests both uppercase and lowercase input formats and ensures that
// unrecognized inputs default to the INFO level.
func TestNewLogLevel(t *testing.T) {
	assert.Equal(t, NewLogLevel("debug"), DEBUG)
	assert.Equal(t, NewLogLevel("DEBUG"), DEBUG)
	assert.Equal(t, NewLogLevel("warn"), WARN)
	assert.Equal(t, NewLogLevel("WARN"), WARN)
	assert.Equal(t, NewLogLevel("error"), ERROR)
	assert.Equal(t, NewLogLevel("ERROR"), ERROR)
	assert.Equal(t, NewLogLevel("panic"), PANIC)
	assert.Equal(t, NewLogLevel("PANIC"), PANIC)
	assert.Equal(t, NewLogLevel("info"), INFO)
}
