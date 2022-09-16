package env

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	RequiredMessagingErrorMessage = "[ConfigBuilder::Messaging] %s is required"
)

func (c *Config) Messaging() ConfigBuilder {
	if c.Err != nil {
		return c
	}

	c.getEngines()
	if c.Err != nil {
		return c
	}

	c.getRabbitMQConfigs()
	if c.Err != nil {
		return c
	}

	c.getKafkaConfigs()
	if c.Err != nil {
		return c
	}

	return c
}

func (c *Config) getEngines() {
	rawEngines := os.Getenv(MESSAGING_ENGINES_ENV_KEY)
	if rawEngines == "" {
		c.Err = fmt.Errorf(RequiredMessagingErrorMessage, MESSAGING_ENGINES_ENV_KEY)
	}

	engSlice := strings.Split(strings.Trim(rawEngines, " "), ",")
	result := map[string]bool{}

	for _, eng := range engSlice {
		switch eng {
		case RABBITMQ_ENGINE:
			result[RABBITMQ_ENGINE] = true
		case KAFKA_ENGINE:
			result[KAFKA_ENGINE] = true
		default:
			c.Err = errors.New("[ConfigBuilder::Messaging] invalid engine")
			return
		}
	}

	c.MESSAGING_ENGINES = result
}

func (c *Config) getRabbitMQConfigs() {
	if _, ok := c.MESSAGING_ENGINES[RABBITMQ_ENGINE]; !ok {
		return
	}

	c.RABBIT_HOST = os.Getenv(RABBIT_HOST_ENV_KEY)
	if c.RABBIT_HOST == "" {
		c.Err = fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_HOST_ENV_KEY)
		return
	}

	c.RABBIT_PORT = os.Getenv(RABBIT_PORT_ENV_KEY)
	if c.RABBIT_PORT == "" {
		c.Err = fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_PORT_ENV_KEY)
		return
	}

	c.RABBIT_USER = os.Getenv(RABBIT_USER_ENV_KEY)
	if c.RABBIT_USER == "" {
		c.Err = fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_USER_ENV_KEY)
		return
	}

	c.RABBIT_PASSWORD = os.Getenv(RABBIT_PASSWORD_ENV_KEY)
	if c.RABBIT_PASSWORD == "" {
		c.Err = fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_PASSWORD_ENV_KEY)
	}

	c.RABBIT_VHOST = os.Getenv(RABBIT_VHOST_ENV_KEY)
	if c.RABBIT_VHOST == "" {
		c.Err = fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_VHOST_ENV_KEY)
	}
}

func (c *Config) getKafkaConfigs() {
	if _, ok := c.MESSAGING_ENGINES[KAFKA_ENGINE]; !ok {
		return
	}
}
