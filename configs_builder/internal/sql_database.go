package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadSQLDatabaseConfigs() (*configs.SQLConfigs, error) {

	sqlConfigs := configs.SQLConfigs{}

	sqlConfigs.Host = os.Getenv(keys.SQLDbHostEnvKey)
	if sqlConfigs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQLDbHostEnvKey)
	}

	sqlConfigs.Port = os.Getenv(keys.SQLDbPortEnvKey)
	if sqlConfigs.Port == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQLDbPortEnvKey)
	}

	sqlConfigs.User = os.Getenv(keys.SQLDbUserEnvKey)
	if sqlConfigs.User == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQLDbUserEnvKey)
	}

	sqlConfigs.Password = os.Getenv(keys.SQLDbPasswordEnvKey)
	if sqlConfigs.Password == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQLDbPasswordEnvKey)
	}

	sqlConfigs.DbName = os.Getenv(keys.SQLDbNameEnvKey)
	if sqlConfigs.DbName == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQLDbNameEnvKey)
	}

	p, err := strconv.Atoi(os.Getenv(keys.SQLDbSecondsToPingEnvKey))
	if err != nil {
		return nil, err
	}

	sqlConfigs.SecondsToPing = p

	return &sqlConfigs, nil
}
