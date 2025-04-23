// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package configs provides a comprehensive configuration framework for GoKit applications.
package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewEnvironment verifies that the NewEnvironment function correctly parses
// various string representations of environments into their corresponding Environment values.
// It tests case sensitivity handling and different naming conventions for each environment.
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
