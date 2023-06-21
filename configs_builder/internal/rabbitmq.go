package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadRabbitMQConfigs() (*configs.RabbitMQConfigs, error) {

	rabbitmqConfigs := configs.RabbitMQConfigs{}

	rabbitmqConfigs.Host = os.Getenv(keys.RABBIT_HOST_ENV_KEY)
	if rabbitmqConfigs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.RABBIT_HOST_ENV_KEY)
	}

	rabbitmqConfigs.Host = os.Getenv(keys.RABBIT_PORT_ENV_KEY)
	if rabbitmqConfigs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.RABBIT_PORT_ENV_KEY)
	}

	rabbitmqConfigs.User = os.Getenv(keys.RABBIT_USER_ENV_KEY)
	if rabbitmqConfigs.User == "" {
		return nil, errors.NewErrRequiredConfig(keys.RABBIT_USER_ENV_KEY)
	}

	rabbitmqConfigs.Password = os.Getenv(keys.RABBIT_PASSWORD_ENV_KEY)
	if rabbitmqConfigs.Password == "" {
		return nil, errors.NewErrRequiredConfig(keys.RABBIT_PASSWORD_ENV_KEY)
	}

	rabbitmqConfigs.VHost = os.Getenv(keys.RABBIT_VHOST_ENV_KEY)

	return &rabbitmqConfigs, nil
}
