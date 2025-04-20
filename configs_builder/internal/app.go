// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadEnvironment() configs.Environment {
	return configs.NewEnvironment(os.Getenv(keys.GoEnvKey))
}

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

func ReadAppName() string {
	name := os.Getenv(keys.AppNameEnvKey)

	if name == "" {
		return keys.DefaultAppName
	}

	return name
}
