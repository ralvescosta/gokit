package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadSqlDatabaseConfigs() (*configs.SqlConfigs, error) {

	sqlConfigs := configs.SqlConfigs{}

	sqlConfigs.Host = os.Getenv(keys.SQL_DB_HOST_ENV_KEY)
	if sqlConfigs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQL_DB_HOST_ENV_KEY)
	}

	sqlConfigs.Port = os.Getenv(keys.SQL_DB_PORT_ENV_KEY)
	if sqlConfigs.Port == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQL_DB_PORT_ENV_KEY)
	}

	sqlConfigs.User = os.Getenv(keys.SQL_DB_USER_ENV_KEY)
	if sqlConfigs.User == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQL_DB_USER_ENV_KEY)
	}

	sqlConfigs.Password = os.Getenv(keys.SQL_DB_PASSWORD_ENV_KEY)
	if sqlConfigs.Password == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQL_DB_PASSWORD_ENV_KEY)
	}

	sqlConfigs.DbName = os.Getenv(keys.SQL_DB_NAME_ENV_KEY)
	if sqlConfigs.DbName == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQL_DB_NAME_ENV_KEY)
	}

	p, err := strconv.Atoi(os.Getenv(keys.SQL_DB_SECONDS_TO_PING_ENV_KEY))
	if err != nil {
		return nil, err
	}

	sqlConfigs.SecondsToPing = p

	return &sqlConfigs, nil
}
