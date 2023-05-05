package configs

type AppConfigs struct {
	GoEnv   Environment
	AppName string

	LogLevel LogLevel
	LogPath  string

	UseSecretManager bool
	SecretKey        string
}
