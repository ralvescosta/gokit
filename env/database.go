package env

import (
	"fmt"
	"os"
	"strconv"
)

const (
	RequiredDatabaseErrorMessage = "[ConfigBuilder::Messaging] %s is required"
)

func (b *ConfigBuilderImpl) SqlDatabase() ConfigBuilder {
	b.sqlDatabase = true
	return b
}

func (c *ConfigBuilderImpl) getSqlDatabaseConfigs() (*SqlConfigs, error) {
	if !c.sqlDatabase {
		return nil, nil
	}

	configs := SqlConfigs{}

	configs.Host = os.Getenv(SQL_DB_HOST_ENV_KEY)
	if configs.Host == "" {
		return nil, fmt.Errorf(RequiredDatabaseErrorMessage, SQL_DB_HOST_ENV_KEY)
	}

	configs.Port = os.Getenv(SQL_DB_PORT_ENV_KEY)
	if configs.Port == "" {
		return nil, fmt.Errorf(RequiredDatabaseErrorMessage, SQL_DB_PORT_ENV_KEY)
	}

	configs.User = os.Getenv(SQL_DB_USER_ENV_KEY)
	if configs.User == "" {
		return nil, fmt.Errorf(RequiredDatabaseErrorMessage, SQL_DB_USER_ENV_KEY)
	}

	configs.Password = os.Getenv(SQL_DB_PASSWORD_ENV_KEY)
	if configs.Password == "" {
		return nil, fmt.Errorf(RequiredDatabaseErrorMessage, SQL_DB_PASSWORD_ENV_KEY)
	}

	configs.DbName = os.Getenv(SQL_DB_NAME_ENV_KEY)
	if configs.DbName == "" {
		return nil, fmt.Errorf(RequiredDatabaseErrorMessage, SQL_DB_NAME_ENV_KEY)
	}

	p, err := strconv.Atoi(os.Getenv(SQL_DB_SECONDS_TO_PING_ENV_KEY))
	if err != nil {
		return nil, err
	}

	configs.SecondsToPing = p

	return &configs, nil
}
