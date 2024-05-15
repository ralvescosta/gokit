package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvironment(t *testing.T) {

	assert.Equal(t, NewEnvironment("development"), DevelopmentEnv)
	assert.Equal(t, NewEnvironment("DEVELOPMENT"), DevelopmentEnv)
	assert.Equal(t, NewEnvironment("dev"), DevelopmentEnv)
	assert.Equal(t, NewEnvironment("production"), ProductionEnv)
	assert.Equal(t, NewEnvironment("PRODUCTION"), ProductionEnv)
	assert.Equal(t, NewEnvironment("prod"), ProductionEnv)
	assert.Equal(t, NewEnvironment("staging"), StagingEnv)
	assert.Equal(t, NewEnvironment("STAGING"), StagingEnv)
	assert.Equal(t, NewEnvironment("stg"), StagingEnv)
	assert.Equal(t, NewEnvironment("qa"), QaEnv)
	assert.Equal(t, NewEnvironment("QA"), QaEnv)
	assert.Equal(t, NewEnvironment("unknown"), UnknownEnv)
}
