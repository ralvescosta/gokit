// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package internal provides helper functions for reading and parsing configuration values
// from environment variables. These functions are used internally by the configs_builder
// package to construct configuration objects.
package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

// ReadEnvironment determines the current runtime environment (development, staging, production)
// by reading from the GO_ENV environment variable
func ReadEnvironment() configs.Environment {
	return configs.NewEnvironment(os.Getenv(keys.GoEnvKey))
}

// ReadAppConfigs constructs a complete AppConfigs object by reading values
// from relevant environment variables
func ReadAppConfigs() *configs.AppConfigs {
	appConfigs := configs.AppConfigs{}

	appConfigs.LogLevel = configs.NewLogLevel(os.Getenv(keys.LogLevelEnvKey))
	appConfigs.AppName = ReadAppName()
	appConfigs.LogPath = os.Getenv(keys.LogPathEnvKey)
	appConfigs.UseSecretManager = func() bool {
		switch os.Getenv(keys.UseSecretManagerEnvKey) {
		case "true":
			return true
		case "false":
			return false
		default:
			return false
		}
	}()
	appConfigs.SecretKey = os.Getenv(keys.SecretKeyEnvKey)

	return &appConfigs
}

// ReadAppName retrieves the application name from environment variables,
// returning a default name if not specified
func ReadAppName() string {
	name := os.Getenv(keys.AppNameEnvKey)

	if name == "" {
		return keys.DefaultAppName
	}

	return name
}
