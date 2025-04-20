// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

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
	ConfigsBuilder interface {
		HTTP() ConfigsBuilder
		Tracing() ConfigsBuilder
		Metrics() ConfigsBuilder
		SQLDatabase() ConfigsBuilder
		Identity() ConfigsBuilder
		MQTT() ConfigsBuilder
		RabbitMQ() ConfigsBuilder
		AWS() ConfigsBuilder
		DynamoDB() ConfigsBuilder
		Build() (*configs.Configs, error)
	}

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

func NewConfigsBuilder() ConfigsBuilder {
	return &configsBuilder{}
}

func (b *configsBuilder) HTTP() ConfigsBuilder {
	b.http = true
	return b
}

func (b *configsBuilder) Tracing() ConfigsBuilder {
	b.tracing = true
	return b
}

func (b *configsBuilder) Metrics() ConfigsBuilder {
	b.metrics = true
	return b
}

func (b *configsBuilder) SQLDatabase() ConfigsBuilder {
	b.sql = true
	return b
}

func (b *configsBuilder) Identity() ConfigsBuilder {
	b.identity = true
	return b
}

func (b *configsBuilder) MQTT() ConfigsBuilder {
	b.mqtt = true
	return b
}

func (b *configsBuilder) RabbitMQ() ConfigsBuilder {
	b.rabbitmq = true
	return b
}

func (b *configsBuilder) AWS() ConfigsBuilder {
	b.aws = true
	return b
}

func (b *configsBuilder) DynamoDB() ConfigsBuilder {
	b.dynamoDB = true
	return b
}

func (b *configsBuilder) Build() (*configs.Configs, error) {
	env := internal.ReadEnvironment()
	if env == configs.UnknownEnv {
		return nil, errors.ErrUnknownEnv
	}

	err := dotEnvConfig(".env." + env.ToString())
	if err != nil {
		return nil, err
	}

	cfgs := configs.Configs{}

	cfgs.AppConfigs = internal.ReadAppConfigs()
	cfgs.AppConfigs.GoEnv = env

	logger, err := logging.NewDefaultLogger(&cfgs)
	if err != nil {
		return nil, err
	}

	cfgs.Logger = logger.(*zap.Logger)

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

var dotEnvConfig = godotenv.Load
