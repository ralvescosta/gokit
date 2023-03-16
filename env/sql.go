package env

import (
	"os"
	"strconv"
)

type (
	SqlConfigs struct {
		Host          string
		Port          string
		User          string
		Password      string
		DbName        string
		SecondsToPing int
	}
)

func (b *ConfigsBuilderImpl) SqlDatabase() ConfigsBuilder {
	b.sqlDatabase = true
	return b
}

func (c *ConfigsBuilderImpl) getSqlDatabaseConfigs() (*SqlConfigs, error) {
	if !c.sqlDatabase {
		return nil, nil
	}

	configs := SqlConfigs{}

	configs.Host = os.Getenv(SQL_DB_HOST_ENV_KEY)
	if configs.Host == "" {
		return nil, NewErrRequiredConfig(SQL_DB_HOST_ENV_KEY)
	}

	configs.Port = os.Getenv(SQL_DB_PORT_ENV_KEY)
	if configs.Port == "" {
		return nil, NewErrRequiredConfig(SQL_DB_PORT_ENV_KEY)
	}

	configs.User = os.Getenv(SQL_DB_USER_ENV_KEY)
	if configs.User == "" {
		return nil, NewErrRequiredConfig(SQL_DB_USER_ENV_KEY)
	}

	configs.Password = os.Getenv(SQL_DB_PASSWORD_ENV_KEY)
	if configs.Password == "" {
		return nil, NewErrRequiredConfig(SQL_DB_PASSWORD_ENV_KEY)
	}

	configs.DbName = os.Getenv(SQL_DB_NAME_ENV_KEY)
	if configs.DbName == "" {
		return nil, NewErrRequiredConfig(SQL_DB_NAME_ENV_KEY)
	}

	p, err := strconv.Atoi(os.Getenv(SQL_DB_SECONDS_TO_PING_ENV_KEY))
	if err != nil {
		return nil, err
	}

	configs.SecondsToPing = p

	return &configs, nil
}
