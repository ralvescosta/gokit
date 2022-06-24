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

func (c *Configs) Messaging() IConfigs {
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

func (c *Configs) getEngines() {
	rawEngines := os.Getenv(MESSAGING_ENGINES_ENV_KEY)
	if rawEngines == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredMessagingErrorMessage, MESSAGING_ENGINES_ENV_KEY))
		return
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

func (c *Configs) getRabbitMQConfigs() {
	if _, ok := c.MESSAGING_ENGINES[RABBITMQ_ENGINE]; !ok {
		return
	}

	c.RABBIT_HOST = os.Getenv(RABBIT_HOST_ENV_KEY)
	if c.RABBIT_HOST == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredMessagingErrorMessage, RABBIT_HOST_ENV_KEY))
		return
	}

	c.RABBIT_PORT = os.Getenv(RABBIT_PORT_ENV_KEY)
	if c.RABBIT_HOST == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredMessagingErrorMessage, RABBIT_PORT_ENV_KEY))
		return
	}

	c.RABBIT_USER = os.Getenv(RABBIT_USER_ENV_KEY)
	if c.RABBIT_HOST == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredMessagingErrorMessage, RABBIT_USER_ENV_KEY))
		return
	}

	c.RABBIT_PASSWORD = os.Getenv(RABBIT_PASSWORD_ENV_KEY)
	if c.RABBIT_HOST == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredMessagingErrorMessage, RABBIT_PASSWORD_ENV_KEY))
	}
}

func (c *Configs) getKafkaConfigs() {
	if _, ok := c.MESSAGING_ENGINES[KAFKA_ENGINE]; !ok {
		return
	}
}
