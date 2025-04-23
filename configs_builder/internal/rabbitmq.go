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

// ReadRabbitMQConfigs retrieves RabbitMQ connection configuration from environment variables.
// Validates that all required connection parameters are provided and returns an error
// if any required configuration is missing. Uses default values for optional parameters.
func ReadRabbitMQConfigs() (*configs.RabbitMQConfigs, error) {
	rabbitmqConfigs := configs.RabbitMQConfigs{}

	// Get connection schema with default fallback to "amqp"
	rabbitmqConfigs.Schema = os.Getenv(keys.RabbitSchemaEnvKey)
	if rabbitmqConfigs.Schema == "" {
		rabbitmqConfigs.Schema = "amqp"
	}

	// Get and validate RabbitMQ host
	rabbitmqConfigs.Host = os.Getenv(keys.RabbitHostEnvKey)
	if rabbitmqConfigs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.RabbitHostEnvKey)
	}

	// Get and validate RabbitMQ port
	rabbitmqConfigs.Port = os.Getenv(keys.RabbitPortEnvKey)
	if rabbitmqConfigs.Port == "" {
		return nil, errors.NewErrRequiredConfig(keys.RabbitPortEnvKey)
	}

	// Get and validate RabbitMQ user
	rabbitmqConfigs.User = os.Getenv(keys.RabbitUserEnvKey)
	if rabbitmqConfigs.User == "" {
		return nil, errors.NewErrRequiredConfig(keys.RabbitUserEnvKey)
	}

	// Get and validate RabbitMQ password
	rabbitmqConfigs.Password = os.Getenv(keys.RabbitPasswordEnvKey)
	if rabbitmqConfigs.Password == "" {
		return nil, errors.NewErrRequiredConfig(keys.RabbitPasswordEnvKey)
	}

	// Get RabbitMQ virtual host (optional)
	rabbitmqConfigs.VHost = os.Getenv(keys.RabbitVHostEnvKey)

	return &rabbitmqConfigs, nil
}
