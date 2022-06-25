package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	RequiredDatabaseErrorMessage = "[ConfigBuilder::Messaging] %s is required"
)

func (c *Configs) Database() IConfigs {
	if c.Err != nil {
		return c
	}

	c.SQL_DB_HOST = os.Getenv(SQL_DB_HOST_ENV_KEY)
	if c.SQL_DB_HOST == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredDatabaseErrorMessage, SQL_DB_HOST_ENV_KEY))
		return c
	}

	c.SQL_DB_PORT = os.Getenv(SQL_DB_PORT_ENV_KEY)
	if c.SQL_DB_PORT == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredDatabaseErrorMessage, SQL_DB_PORT_ENV_KEY))
		return c
	}

	c.SQL_DB_USER = os.Getenv(SQL_DB_USER_ENV_KEY)
	if c.SQL_DB_USER == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredDatabaseErrorMessage, SQL_DB_USER_ENV_KEY))
		return c
	}

	c.SQL_DB_PASSWORD = os.Getenv(SQL_DB_PASSWORD_ENV_KEY)
	if c.SQL_DB_PASSWORD == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredDatabaseErrorMessage, SQL_DB_PASSWORD_ENV_KEY))
		return c
	}

	c.SQL_DB_NAME = os.Getenv(SQL_DB_NAME_ENV_KEY)
	if c.SQL_DB_NAME == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredDatabaseErrorMessage, SQL_DB_NAME_ENV_KEY))
		return c
	}

	p, err := strconv.Atoi(os.Getenv(SQL_DB_SECONDS_TO_PING_ENV_KEY))
	if err != nil {
		c.Err = err
		return c
	}

	c.SQL_DB_SECONDS_TO_PING = p

	return c
}
