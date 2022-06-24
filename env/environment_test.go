package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvironment(t *testing.T) {

	assert.Equal(t, NewEnvironment("development"), DEVELOPMENT_ENV)
	assert.Equal(t, NewEnvironment("DEVELOPMENT"), DEVELOPMENT_ENV)
	assert.Equal(t, NewEnvironment("dev"), DEVELOPMENT_ENV)
	assert.Equal(t, NewEnvironment("production"), PRODUCTION_ENV)
	assert.Equal(t, NewEnvironment("PRODUCTION"), PRODUCTION_ENV)
	assert.Equal(t, NewEnvironment("prod"), PRODUCTION_ENV)
	assert.Equal(t, NewEnvironment("staging"), STAGING_ENV)
	assert.Equal(t, NewEnvironment("STAGING"), STAGING_ENV)
	assert.Equal(t, NewEnvironment("stg"), STAGING_ENV)
	assert.Equal(t, NewEnvironment("qa"), QA_ENV)
	assert.Equal(t, NewEnvironment("QA"), QA_ENV)
	assert.Equal(t, NewEnvironment("unknown"), UNKNOWN_ENV)
}
