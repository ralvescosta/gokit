package env

import (
	"os"
)

type (
	RabbitMQConfigs struct {
		Host     string
		Port     string
		User     string
		Password string
		VHost    string
	}
)

func (b *ConfigsBuilderImpl) RabbitMQ() ConfigsBuilder {
	b.rabbitmq = true
	return b
}

func (b *ConfigsBuilderImpl) getRabbitMQConfigs() (*RabbitMQConfigs, error) {
	if !b.rabbitmq {
		return nil, nil
	}

	configs := RabbitMQConfigs{}

	configs.Host = os.Getenv(RABBIT_HOST_ENV_KEY)
	if configs.Host == "" {
		return nil, NewErrRequiredConfig(RABBIT_HOST_ENV_KEY)
	}

	configs.Host = os.Getenv(RABBIT_PORT_ENV_KEY)
	if configs.Host == "" {
		return nil, NewErrRequiredConfig(RABBIT_PORT_ENV_KEY)
	}

	configs.User = os.Getenv(RABBIT_USER_ENV_KEY)
	if configs.User == "" {
		return nil, NewErrRequiredConfig(RABBIT_USER_ENV_KEY)
	}

	configs.Password = os.Getenv(RABBIT_PASSWORD_ENV_KEY)
	if configs.Password == "" {
		return nil, NewErrRequiredConfig(RABBIT_PASSWORD_ENV_KEY)
	}

	configs.VHost = os.Getenv(RABBIT_VHOST_ENV_KEY)

	return &configs, nil
}
