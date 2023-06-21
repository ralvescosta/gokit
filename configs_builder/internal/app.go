package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadEnvironment() configs.Environment {
	return configs.NewEnvironment(os.Getenv(keys.GO_ENV_KEY))
}

func ReadAppConfigs() *configs.AppConfigs {
	appConfigs := configs.AppConfigs{}

	appConfigs.LogLevel = configs.NewLogLevel(os.Getenv(keys.LOG_LEVEL_ENV_KEY))
	appConfigs.AppName = ReadAppName()
	appConfigs.LogPath = os.Getenv(keys.LOG_PATH_ENV_KEY)
	appConfigs.UseSecretManager = func() bool {
		switch os.Getenv(keys.USE_SECRET_MANAGER_ENV_KEY) {
		case "true":
			return true
		case "false":
			return false
		default:
			return false
		}
	}()
	appConfigs.SecretKey = os.Getenv(keys.SECRET_KEY_ENV_KEY)

	return &appConfigs
}

func ReadAppName() string {
	name := os.Getenv(keys.APP_NAME_ENV_KEY)

	if name == "" {
		return keys.DEFAULT_APP_NAME
	}

	return name
}
