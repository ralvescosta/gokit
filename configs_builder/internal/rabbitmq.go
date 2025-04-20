// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadRabbitMQConfigs() (*configs.RabbitMQConfigs, error) {

	rabbitmqConfigs := configs.RabbitMQConfigs{}

	rabbitmqConfigs.Host = os.Getenv(keys.RabbitHostEnvKey)
	if rabbitmqConfigs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.RabbitHostEnvKey)
	}

	rabbitmqConfigs.Port = os.Getenv(keys.RabbitPortEnvKey)
	if rabbitmqConfigs.Port == "" {
		return nil, errors.NewErrRequiredConfig(keys.RabbitPortEnvKey)
	}

	rabbitmqConfigs.User = os.Getenv(keys.RabbitUserEnvKey)
	if rabbitmqConfigs.User == "" {
		return nil, errors.NewErrRequiredConfig(keys.RabbitUserEnvKey)
	}

	rabbitmqConfigs.Password = os.Getenv(keys.RabbitPasswordEnvKey)
	if rabbitmqConfigs.Password == "" {
		return nil, errors.NewErrRequiredConfig(keys.RabbitPasswordEnvKey)
	}

	rabbitmqConfigs.VHost = os.Getenv(keys.RabbitVHostEnvKey)

	return &rabbitmqConfigs, nil
}
