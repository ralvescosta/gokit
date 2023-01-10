package env

import (
	"fmt"
	"os"
)

const (
	RequiredMessagingErrorMessage = "[ConfigBuilder::Messaging] %s is required"
)

func (c *Config) RabbitMQ() ConfigBuilder {
	if c.Err != nil {
		return c
	}

	c.RABBIT_HOST = os.Getenv(RABBIT_HOST_ENV_KEY)
	if c.RABBIT_HOST == "" {
		c.Err = fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_HOST_ENV_KEY)
		return c
	}

	c.RABBIT_PORT = os.Getenv(RABBIT_PORT_ENV_KEY)
	if c.RABBIT_PORT == "" {
		c.Err = fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_PORT_ENV_KEY)
		return c
	}

	c.RABBIT_USER = os.Getenv(RABBIT_USER_ENV_KEY)
	if c.RABBIT_USER == "" {
		c.Err = fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_USER_ENV_KEY)
		return c
	}

	c.RABBIT_PASSWORD = os.Getenv(RABBIT_PASSWORD_ENV_KEY)
	if c.RABBIT_PASSWORD == "" {
		c.Err = fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_PASSWORD_ENV_KEY)
		return c
	}

	c.RABBIT_VHOST = os.Getenv(RABBIT_VHOST_ENV_KEY)

	return c
}
