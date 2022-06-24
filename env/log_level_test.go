package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogLevel(t *testing.T) {
	assert.Equal(t, NewLogLevel("debug"), DEBUG_L)
	assert.Equal(t, NewLogLevel("DEBUG"), DEBUG_L)
	assert.Equal(t, NewLogLevel("warn"), WARN_L)
	assert.Equal(t, NewLogLevel("WARN"), WARN_L)
	assert.Equal(t, NewLogLevel("error"), ERROR_L)
	assert.Equal(t, NewLogLevel("ERROR"), ERROR_L)
	assert.Equal(t, NewLogLevel("panic"), PANIC_L)
	assert.Equal(t, NewLogLevel("PANIC"), PANIC_L)
	assert.Equal(t, NewLogLevel("info"), INFO_L)
}
