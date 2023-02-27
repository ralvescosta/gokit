package env

import (
	"fmt"
	"os"
)

const (
	RequiredMessagingErrorMessage = "[ConfigBuilder::Messaging] %s is required"
)

func (b *ConfigBuilderImpl) RabbitMQ() ConfigBuilder {
	b.rabbitmq = true
	return b
}

func (b *ConfigBuilderImpl) getRabbitMQConfigs() (*RabbitMQConfigs, error) {
	if !b.rabbitmq {
		return nil, nil
	}

	configs := RabbitMQConfigs{}

	configs.Host = os.Getenv(RABBIT_HOST_ENV_KEY)
	if configs.Host == "" {
		return nil, fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_HOST_ENV_KEY)
	}

	configs.Host = os.Getenv(RABBIT_PORT_ENV_KEY)
	if configs.Host == "" {
		return nil, fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_PORT_ENV_KEY)
	}

	configs.User = os.Getenv(RABBIT_USER_ENV_KEY)
	if configs.User == "" {
		return nil, fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_USER_ENV_KEY)
	}

	configs.Password = os.Getenv(RABBIT_PASSWORD_ENV_KEY)
	if configs.Password == "" {
		return nil, fmt.Errorf(RequiredMessagingErrorMessage, RABBIT_PASSWORD_ENV_KEY)
	}

	configs.VHost = os.Getenv(RABBIT_VHOST_ENV_KEY)

	return &configs, nil
}
