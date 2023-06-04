package configsbuilder

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ralvescosta/dotenv"
	"github.com/ralvescosta/gokit/configs"
)

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
		Build() (interface{}, error)
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
		dynamoDB    bool
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

func (b *configsBuilder) Build() (interface{}, error) {
	appConfigs, err := b.readAppConfigs()
	if err != nil {
		return nil, err
	}

	httpServerConfigs, err := b.readHTTPConfigs()
	if err != nil {
		return nil, err
	}

	otelConfigs, err := b.readOtelConfigs()
	if err != nil {
		return nil, err
	}

	sqlDatabaseConfigs, err := b.readSqlDatabaseConfigs()
	if err != nil {
		return nil, err
	}

	auth0Configs, err := b.readAuth0Configs()
	if err != nil {
		return nil, err
	}

	mqttConfigs, err := b.readMQTTConfigs()
	if err != nil {
		return nil, err
	}

	rabbitMQConfigs, err := b.readRabbitMQConfigs()
	if err != nil {
		return nil, err
	}

	awsConfigs, err := b.readAWSConfigs()
	if err != nil {
		return nil, err
	}

	dynamoDbConfigs, err := b.readDynamoDBConfigs()
	if err != nil {
		return nil, err
	}

	return &configs.Configs{
		AppConfigs:      appConfigs,
		HTTPConfigs:     httpServerConfigs,
		OtelConfigs:     otelConfigs,
		SqlConfigs:      sqlDatabaseConfigs,
		Auth0Configs:    auth0Configs,
		MQTTConfigs:     mqttConfigs,
		RabbitMQConfigs: rabbitMQConfigs,
		AWSConfigs:      awsConfigs,
		DynamoDBConfigs: dynamoDbConfigs,
	}, nil
}

func (b *configsBuilder) readAppConfigs() (*configs.AppConfigs, error) {
	appConfigs := configs.AppConfigs{}
	appConfigs.GoEnv = configs.NewEnvironment(os.Getenv(GO_ENV_KEY))

	if appConfigs.GoEnv == configs.UNKNOWN_ENV {
		return nil, ErrUnknownEnv
	}

	err := dotEnvConfig(".env." + appConfigs.GoEnv.ToString())
	if err != nil {
		return nil, err
	}

	appConfigs.LogLevel = configs.NewLogLevel(os.Getenv(LOG_LEVEL_ENV_KEY))
	appConfigs.AppName = b.appName()

	return &appConfigs, nil
}

func (b *configsBuilder) appName() string {
	name := os.Getenv(APP_NAME_ENV_KEY)

	if name == "" {
		return DEFAULT_APP_NAME
	}

	return name
}

func (b *configsBuilder) readHTTPConfigs() (*configs.HTTPConfigs, error) {
	if !b.http {
		return nil, nil
	}

	httpConfigs := configs.HTTPConfigs{}

	httpConfigs.Port = os.Getenv(HTTP_PORT_ENV_KEY)
	if httpConfigs.Port == "" {
		return nil, NewErrRequiredConfig(HTTP_PORT_ENV_KEY)
	}

	httpConfigs.Host = os.Getenv(HTTP_HOST_ENV_KEY)
	if httpConfigs.Host == "" {
		return nil, NewErrRequiredConfig(HTTP_HOST_ENV_KEY)
	}

	httpConfigs.Addr = fmt.Sprintf("%s:%s", httpConfigs.Host, httpConfigs.Port)

	profiling := os.Getenv(HTTP_ENABLE_PROFILING_ENV_KEY)
	if profiling == "true" {
		httpConfigs.EnableProfiling = true
	}

	return &httpConfigs, nil
}

func (b *configsBuilder) readOtelConfigs() (*configs.OtelConfigs, error) {
	if !b.otel {
		return nil, nil
	}

	tracingEnabled := os.Getenv(TRACING_ENABLED_ENV_KEY)
	metricsEnabled := os.Getenv(METRICS_ENABLED_ENV_KEY)

	if tracingEnabled == "" && metricsEnabled == "" {
		return nil, NewErrRequiredConfig(TRACING_ENABLED_ENV_KEY)
	}

	otelConfigs := configs.OtelConfigs{}

	if tracingEnabled == "true" {
		otelConfigs.TracingEnabled = true
	}

	if metricsEnabled == "true" {
		otelConfigs.MetricsEnabled = true
	}

	otelConfigs.OtlpEndpoint = os.Getenv(OTLP_ENDPOINT_ENV_KEY)
	otelConfigs.OtlpApiKey = os.Getenv(OTLP_API_KEY_ENV_KEY)
	otelConfigs.JaegerServiceName = os.Getenv(JAEGER_SERVICE_NAME_KEY)
	otelConfigs.JaegerAgentHost = os.Getenv(JAEGER_AGENT_HOST_KEY)
	otelConfigs.JaegerSampleType = os.Getenv(JAEGER_SAMPLER_TYPE_KEY)
	if samplerParam := os.Getenv(JAEGER_SAMPLER_PARAM_KEY); samplerParam != "" {
		otelConfigs.JaegerSampleParam, _ = strconv.Atoi(samplerParam)
	}

	if reportLogSpans := os.Getenv(JAEGER_REPORTER_LOG_SPANS_KEY); reportLogSpans != "" {
		otelConfigs.JaegerReporterLogSpans = reportLogSpans == "true"
	}

	if rpcMetrics := os.Getenv(JAEGER_RPC_METRICS_KEY); rpcMetrics != "" {
		otelConfigs.JaegerRpcMetrics = rpcMetrics == "true"
	}

	return &otelConfigs, nil
}

func (b *configsBuilder) readSqlDatabaseConfigs() (*configs.SqlConfigs, error) {
	if !b.sqlDatabase {
		return nil, nil
	}

	sqlConfigs := configs.SqlConfigs{}

	sqlConfigs.Host = os.Getenv(SQL_DB_HOST_ENV_KEY)
	if sqlConfigs.Host == "" {
		return nil, NewErrRequiredConfig(SQL_DB_HOST_ENV_KEY)
	}

	sqlConfigs.Port = os.Getenv(SQL_DB_PORT_ENV_KEY)
	if sqlConfigs.Port == "" {
		return nil, NewErrRequiredConfig(SQL_DB_PORT_ENV_KEY)
	}

	sqlConfigs.User = os.Getenv(SQL_DB_USER_ENV_KEY)
	if sqlConfigs.User == "" {
		return nil, NewErrRequiredConfig(SQL_DB_USER_ENV_KEY)
	}

	sqlConfigs.Password = os.Getenv(SQL_DB_PASSWORD_ENV_KEY)
	if sqlConfigs.Password == "" {
		return nil, NewErrRequiredConfig(SQL_DB_PASSWORD_ENV_KEY)
	}

	sqlConfigs.DbName = os.Getenv(SQL_DB_NAME_ENV_KEY)
	if sqlConfigs.DbName == "" {
		return nil, NewErrRequiredConfig(SQL_DB_NAME_ENV_KEY)
	}

	p, err := strconv.Atoi(os.Getenv(SQL_DB_SECONDS_TO_PING_ENV_KEY))
	if err != nil {
		return nil, err
	}

	sqlConfigs.SecondsToPing = p

	return &sqlConfigs, nil
}

func (b *configsBuilder) readAuth0Configs() (*configs.Auth0Configs, error) {
	return nil, nil
}

func (b *configsBuilder) readMQTTConfigs() (*configs.MQTTConfigs, error) {
	return nil, nil
}

func (b *configsBuilder) readRabbitMQConfigs() (*configs.RabbitMQConfigs, error) {
	if !b.rabbitmq {
		return nil, nil
	}

	rabbitmqConfigs := configs.RabbitMQConfigs{}

	rabbitmqConfigs.Host = os.Getenv(RABBIT_HOST_ENV_KEY)
	if rabbitmqConfigs.Host == "" {
		return nil, NewErrRequiredConfig(RABBIT_HOST_ENV_KEY)
	}

	rabbitmqConfigs.Host = os.Getenv(RABBIT_PORT_ENV_KEY)
	if rabbitmqConfigs.Host == "" {
		return nil, NewErrRequiredConfig(RABBIT_PORT_ENV_KEY)
	}

	rabbitmqConfigs.User = os.Getenv(RABBIT_USER_ENV_KEY)
	if rabbitmqConfigs.User == "" {
		return nil, NewErrRequiredConfig(RABBIT_USER_ENV_KEY)
	}

	rabbitmqConfigs.Password = os.Getenv(RABBIT_PASSWORD_ENV_KEY)
	if rabbitmqConfigs.Password == "" {
		return nil, NewErrRequiredConfig(RABBIT_PASSWORD_ENV_KEY)
	}

	rabbitmqConfigs.VHost = os.Getenv(RABBIT_VHOST_ENV_KEY)

	return &rabbitmqConfigs, nil
}

func (b *configsBuilder) readAWSConfigs() (*configs.AWSConfigs, error) {
	return nil, nil
}

func (b *configsBuilder) readDynamoDBConfigs() (*configs.DynamoDBConfigs, error) {
	return nil, nil
}

var dotEnvConfig = dotenv.Configure
