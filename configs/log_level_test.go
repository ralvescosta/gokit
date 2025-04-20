// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
