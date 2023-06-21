package configsbuilder

import (
	"github.com/ralvescosta/dotenv"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	"github.com/ralvescosta/gokit/configs_builder/internal"
)

type (
	ConfigsBuilder interface {
		HTTP() ConfigsBuilder
		Otel() ConfigsBuilder
		PostgreSQL() ConfigsBuilder
		Auth0() ConfigsBuilder
		MQTT() ConfigsBuilder
		RabbitMQ() ConfigsBuilder
		AWS() ConfigsBuilder
		DynamoDB() ConfigsBuilder
		Build() (interface{}, error)
	}

	configsBuilder struct {
		Err error

		http     bool
		otel     bool
		postgres bool
		auth0    bool
		mqtt     bool
		rabbitmq bool
		aws      bool
		dynamoDB bool
	}
)

func NewConfigsBuilder() *configsBuilder {
	return &configsBuilder{}
}

func (b *configsBuilder) HTTP() ConfigsBuilder {
	b.http = true
	return b
}

func (b *configsBuilder) Otel() ConfigsBuilder {
	b.otel = true
	return b
}

func (b *configsBuilder) PostgreSQL() ConfigsBuilder {
	b.postgres = true
	return b
}

func (b *configsBuilder) Auth0() ConfigsBuilder {
	b.auth0 = true
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

func (b *configsBuilder) Build() (interface{}, error) {
	env := internal.ReadEnvironment()
	if env == configs.UNKNOWN_ENV {
		return nil, errors.ErrUnknownEnv
	}

	err := dotEnvConfig(".env." + env.ToString())
	if err != nil {
		return nil, err
	}

	cfgs := configs.Configs{}

	cfgs.AppConfigs = internal.ReadAppConfigs()
	cfgs.AppConfigs.GoEnv = env

	cfgs.HTTPConfigs, err = internal.ReadHTTPConfigs()
	if err != nil {
		return nil, err
	}

	cfgs.OtelConfigs, err = internal.ReadOtelConfigs()
	if err != nil {
		return nil, err
	}

	cfgs.SqlConfigs, err = internal.ReadSqlDatabaseConfigs()
	if err != nil {
		return nil, err
	}

	cfgs.Auth0Configs, err = internal.ReadAuth0Configs()
	if err != nil {
		return nil, err
	}

	cfgs.MQTTConfigs, err = internal.ReadMQTTConfigs()
	if err != nil {
		return nil, err
	}

	cfgs.RabbitMQConfigs, err = internal.ReadRabbitMQConfigs()
	if err != nil {
		return nil, err
	}

	cfgs.AWSConfigs, err = internal.ReadAWSConfigs()
	if err != nil {
		return nil, err
	}

	cfgs.DynamoDBConfigs, err = internal.ReadDynamoDBConfigs()
	if err != nil {
		return nil, err
	}

	return &cfgs, nil
}

var dotEnvConfig = dotenv.Configure
