// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// AppConfigs contains the core application configuration settings.
// It defines essential parameters that govern the application's behavior
// such as environment, logging settings, and secret management options.
type AppConfigs struct {
	// GoEnv specifies the environment in which the application runs (e.g., development, production)
	GoEnv Environment
	// AppName holds the name of the application for identification in logs and monitoring
	AppName string

	// LogLevel determines the verbosity and severity threshold of application logs
	LogLevel LogLevel
	// LogPath specifies the file path where logs should be written (if file logging is enabled)
	LogPath string

	// UseSecretManager indicates whether the application should retrieve secrets from a secure secret manager
	UseSecretManager bool
	// SecretKey identifies the key/path used to fetch secrets from the secret manager
	SecretKey string
}
