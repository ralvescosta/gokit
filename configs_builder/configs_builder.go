package configsbuilder

import "os"

const (
	GO_ENV_KEY        = "GO_ENV"
	LOG_LEVEL_ENV_KEY = "LOG_LEVEL"
	LOG_PATH_ENV_KEY  = "LOG_PATH"
	APP_NAME_ENV_KEY  = "APP_NAME"

	SQL_DB_HOST_ENV_KEY            = "SQL_DB_HOST"
	SQL_DB_PORT_ENV_KEY            = "SQL_DB_PORT"
	SQL_DB_USER_ENV_KEY            = "SQL_DB_USER"
	SQL_DB_PASSWORD_ENV_KEY        = "SQL_DB_PASSWORD"
	SQL_DB_NAME_ENV_KEY            = "SQL_DB_NAME"
	SQL_DB_SECONDS_TO_PING_ENV_KEY = "SQL_DB_SECONDS_TO_PING"

	RABBIT_HOST_ENV_KEY     = "RABBIT_HOST"
	RABBIT_PORT_ENV_KEY     = "RABBIT_PORT"
	RABBIT_USER_ENV_KEY     = "RABBIT_USER"
	RABBIT_PASSWORD_ENV_KEY = "RABBIT_PASSWORD"
	RABBIT_VHOST_ENV_KEY    = "RABBIT_VHOST"
	KAFKA_HOST_ENV_KEY      = "KAFKA_HOST"
	KAFKA_PORT_ENV_KEY      = "KAFKA_PORT"
	KAFKA_USER_ENV_KEY      = "KAFKA_USER"
	KAFKA_PASSWORD_ENV_KEY  = "KAFKA_PASSWORD"
	RABBITMQ_ENGINE         = "RabbitMQ"
	KAFKA_ENGINE            = "Kafka"

	DEFAULT_APP_NAME = "app"
	DEFAULT_LOG_PATH = "/logs/"

	TRACING_ENABLED_ENV_KEY       = "TRACING_ENABLED"
	METRICS_ENABLED_ENV_KEY       = "METRICS_ENABLE"
	OTLP_ENDPOINT_ENV_KEY         = "OTLP_ENDPOINT"
	OTLP_API_KEY_ENV_KEY          = "OTLP_API_KEY"
	JAEGER_SERVICE_NAME_KEY       = "JAEGER_SERVICE_NAME"
	JAEGER_AGENT_HOST_KEY         = "JAEGER_AGENT_HOST"
	JAEGER_SAMPLER_TYPE_KEY       = "JAEGER_SAMPLER_TYPE"
	JAEGER_SAMPLER_PARAM_KEY      = "JAEGER_SAMPLER_PARAM"
	JAEGER_REPORTER_LOG_SPANS_KEY = "JAEGER_REPORTER_LOG_SPANS"
	JAEGER_RPC_METRICS_KEY        = "JAEGER_RPC_METRICS"

	HTTP_PORT_ENV_KEY             = "HTTP_PORT"
	HTTP_HOST_ENV_KEY             = "HTTP_HOST"
	HTTP_ENABLE_PROFILING_ENV_KEY = "HTTP_ENABLE_PROFILING"
)

type (
	ConfigsBuilder interface {
		HTTP() ConfigsBuilder
		Otel() ConfigsBuilder
		SqlDatabase() ConfigsBuilder
		Auth0() ConfigsBuilder
		MQTT() ConfigsBuilder
		RabbitMQ() ConfigsBuilder
		AWS() ConfigsBuilder
		DynamoDB() ConfigsBuilder
		Build() (*Configs, error)
	}

	configsBuilder struct {
		Err error

		http        bool
		otel        bool
		sqlDatabase bool
		auth0       bool
		mqtt        bool
		rabbitmq    bool
		aws         bool
		dynamoDB
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

func (b *configsBuilder) SqlDatabase() ConfigsBuilder {
	b.sqlDatabase = true
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

func (b *configsBuilder) Build() (*Configs, error) {
	// appConfigs, err := b.getAppConfigs()
	// if err != nil {
	// 	return nil, err
	// }

	// sqlDatabaseConfigs, err := b.getSqlDatabaseConfigs()
	// if err != nil {
	// 	return nil, err
	// }

	// rabbitMQConfigs, err := b.getRabbitMQConfigs()
	// if err != nil {
	// 	return nil, err
	// }

	// otelConfigs, err := b.getOtelConfigs()
	// if err != nil {
	// 	return nil, err
	// }

	// httpServerConfigs, err := b.getHTTPServerConfigs()
	// if err != nil {
	// 	return nil, err
	// }

	// return &Configs{
	// 	AppConfigs:      appConfigs,
	// 	SqlConfigs:      sqlDatabaseConfigs,
	// 	RabbitMQConfigs: rabbitMQConfigs,
	// 	OtelConfigs:     otelConfigs,
	// 	HTTPConfigs:     httpServerConfigs,
	// }, nil
}

func (b *ConfigsBuilderImpl) getAppConfigs() (*AppConfigs, error) {
	configs := AppConfigs{}
	configs.GoEnv = NewEnvironment(os.Getenv(GO_ENV_KEY))

	if configs.GoEnv == UNKNOWN_ENV {
		return nil, ErrUnknownEnv
	}

	err := dotEnvConfig(".env." + configs.GoEnv.ToString())
	if err != nil {
		return nil, err
	}

	configs.LogLevel = NewLogLevel(os.Getenv(LOG_LEVEL_ENV_KEY))
	configs.AppName = b.appName()

	return &configs, nil
}

func (b *ConfigsBuilderImpl) appName() string {
	name := os.Getenv(APP_NAME_ENV_KEY)

	if name == "" {
		return DEFAULT_APP_NAME
	}

	return name
}
