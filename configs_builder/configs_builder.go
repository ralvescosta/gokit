// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package configsbuilder provides a fluent interface for building application configurations.
// It simplifies the process of loading configurations from environment variables and .env files
// for various components of an application such as HTTP, messaging, databases, etc.
package configsbuilder

import (
	"github.com/joho/godotenv"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"

	"github.com/ralvescosta/gokit/configs_builder/errors"
	"github.com/ralvescosta/gokit/configs_builder/internal"
)

type (
	// ConfigsBuilder defines the interface for the builder pattern used to construct configurations.
	// It provides methods to specify which configuration components should be included.
	ConfigsBuilder interface {
		// HTTP enables HTTP server configuration loading
		HTTP() ConfigsBuilder
		// Tracing enables OpenTelemetry tracing configuration loading
		Tracing() ConfigsBuilder
		// Metrics enables metrics configuration loading
		Metrics() ConfigsBuilder
		// SQLDatabase enables SQL database configuration loading
		SQLDatabase() ConfigsBuilder
		// Identity enables identity/authentication configuration loading
		Identity() ConfigsBuilder
		// MQTT enables MQTT client configuration loading
		MQTT() ConfigsBuilder
		// RabbitMQ enables RabbitMQ configuration loading
		RabbitMQ() ConfigsBuilder
		// AWS enables AWS configuration loading
		AWS() ConfigsBuilder
		// DynamoDB enables DynamoDB configuration loading
		DynamoDB() ConfigsBuilder
		// Build processes all enabled configurations and returns the complete config object
		Build() (*configs.Configs, error)
	}

	// configsBuilder implements the ConfigsBuilder interface and tracks which configurations to load
	configsBuilder struct {
		Err error

		http     bool
		tracing  bool
		metrics  bool
		sql      bool
		identity bool
		mqtt     bool
		rabbitmq bool
		aws      bool
		dynamoDB bool
	}
)

// NewConfigsBuilder creates a new instance of ConfigsBuilder with no configurations enabled
func NewConfigsBuilder() ConfigsBuilder {
	return &configsBuilder{}
}

// HTTP enables HTTP configuration loading in the builder
func (b *configsBuilder) HTTP() ConfigsBuilder {
	b.http = true
	return b
}

// Tracing enables tracing configuration loading in the builder
func (b *configsBuilder) Tracing() ConfigsBuilder {
	b.tracing = true
	return b
}

// Metrics enables metrics configuration loading in the builder
func (b *configsBuilder) Metrics() ConfigsBuilder {
	b.metrics = true
	return b
}

// SQLDatabase enables SQL database configuration loading in the builder
func (b *configsBuilder) SQLDatabase() ConfigsBuilder {
	b.sql = true
	return b
}

// Identity enables identity configuration loading in the builder
func (b *configsBuilder) Identity() ConfigsBuilder {
	b.identity = true
	return b
}

// MQTT enables MQTT configuration loading in the builder
func (b *configsBuilder) MQTT() ConfigsBuilder {
	b.mqtt = true
	return b
}

// RabbitMQ enables RabbitMQ configuration loading in the builder
func (b *configsBuilder) RabbitMQ() ConfigsBuilder {
	b.rabbitmq = true
	return b
}

// AWS enables AWS configuration loading in the builder
func (b *configsBuilder) AWS() ConfigsBuilder {
	b.aws = true
	return b
}

// DynamoDB enables DynamoDB configuration loading in the builder
func (b *configsBuilder) DynamoDB() ConfigsBuilder {
	b.dynamoDB = true
	return b
}

// Build processes all enabled configurations and returns the complete configs object.
// It reads environment variables, loads .env files, and constructs the configuration
// based on the enabled features. Returns an error if any configuration fails to load.
func (b *configsBuilder) Build() (*configs.Configs, error) {
	// Determine the runtime environment
	env := internal.ReadEnvironment()
	if env == configs.UnknownEnv {
		return nil, errors.ErrUnknownEnv
	}

	// Load environment-specific .env file
	err := dotEnvConfig(".env." + env.ToString())
	if err != nil {
		return nil, err
	}

	cfgs := configs.Configs{}

	// Load application base configurations
	cfgs.AppConfigs = internal.ReadAppConfigs()
	cfgs.AppConfigs.GoEnv = env

	// Initialize the logger
	logger, err := logging.NewDefaultLogger(&cfgs)
	if err != nil {
		return nil, err
	}

	cfgs.Logger = logger.(*zap.Logger)

	// Load component-specific configurations based on what was enabled
	if b.http {
		cfgs.HTTPConfigs, err = internal.ReadHTTPConfigs()
		if err != nil {
			return nil, err
		}
	}

	if b.metrics {
		cfgs.MetricsConfigs, err = internal.ReadMetricsConfigs()
		if err != nil {
			return nil, err
		}
	}

	if b.tracing {
		cfgs.TracingConfigs, err = internal.ReadTracingConfigs()
		if err != nil {
			return nil, err
		}
	}

	if b.sql {
		cfgs.SQLConfigs, err = internal.ReadSQLDatabaseConfigs()
		if err != nil {
			return nil, err
		}
	}

	if b.identity {
		cfgs.IdentityConfigs, err = internal.ReadIdentityConfigs()
		if err != nil {
			return nil, err
		}
	}

	if b.mqtt {
		cfgs.MQTTConfigs, err = internal.ReadMQTTConfigs()
		if err != nil {
			return nil, err
		}
	}

	if b.rabbitmq {
		cfgs.RabbitMQConfigs, err = internal.ReadRabbitMQConfigs()
		if err != nil {
			return nil, err
		}
	}

	if b.aws {
		cfgs.AWSConfigs, err = internal.ReadAWSConfigs()
		if err != nil {
			return nil, err
		}
	}

	if b.dynamoDB {
		cfgs.DynamoDBConfigs, err = internal.ReadDynamoDBConfigs()
		if err != nil {
			return nil, err
		}
	}

	return &cfgs, nil
}

// dotEnvConfig is a variable containing the godotenv.Load function, which allows for testing
var dotEnvConfig = godotenv.Load
