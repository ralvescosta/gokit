// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

type AppConfigs struct {
	GoEnv   Environment
	AppName string

	LogLevel LogLevel
	LogPath  string

	UseSecretManager bool
	SecretKey        string
}
