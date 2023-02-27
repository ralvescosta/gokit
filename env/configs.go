package env

import (
	"errors"
	"fmt"
	"os"

	"github.com/ralvescosta/dotenv"
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

	UNKNOWN_ENV     Environment = 0
	LOCAL_ENV       Environment = 1
	DEVELOPMENT_ENV Environment = 2
	STAGING_ENV     Environment = 3
	QA_ENV          Environment = 4
	PRODUCTION_ENV  Environment = 5

	DEBUG_L LogLevel = 0
	INFO_L  LogLevel = 1
	WARN_L  LogLevel = 2
	ERROR_L LogLevel = 3
	PANIC_L LogLevel = 4

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

	HTTP_PORT_ENV_KEY = "HTTP_PORT"
	HTTP_HOST_ENV_KEY = "HTTP_HOST"
)

var (
	EnvironmentMapping = map[Environment]string{
		UNKNOWN_ENV:     "unknown",
		LOCAL_ENV:       "local",
		DEVELOPMENT_ENV: "development",
		STAGING_ENV:     "staging",
		QA_ENV:          "qa",
		PRODUCTION_ENV:  "production",
	}
)

type (
	ConfigBuilder interface {
		SqlDatabase() ConfigBuilder
		RabbitMQ() ConfigBuilder
		Otel() ConfigBuilder
		HTTPServer() ConfigBuilder
		Build() (*Configs, error)
	}

	ConfigBuilderImpl struct {
		Err error

		sqlDatabase bool
		rabbitmq    bool
		otel        bool
		httpServer  bool
	}

	AppConfigs struct {
		GoEnv     Environment
		LogLevel  LogLevel
		LogPath   string
		AppName   string
		SecretKey string
	}

	SqlConfigs struct {
		Host          string
		Port          string
		User          string
		Password      string
		DbName        string
		SecondsToPing int
	}

	RabbitMQConfigs struct {
		Host     string
		Port     string
		User     string
		Password string
		VHost    string
	}

	OtelConfigs struct {
		TracingEnabled         bool
		MetricsEnabled         bool
		OtlpEndpoint           string
		OtlpApiKey             string
		JaegerServiceName      string
		JaegerAgentHost        string
		JaegerSampleType       string
		JaegerSampleParam      int
		JaegerReporterLogSpans bool
		JaegerRpcMetrics       bool
	}

	HTTPConfigs struct {
		Host string
		Port string
		Addr string
	}

	Configs struct {
		Custom map[string]string

		AppConfigs      *AppConfigs
		SqlConfigs      *SqlConfigs
		RabbitMqConfigs *RabbitMQConfigs
		OtelConfigs     *OtelConfigs
		HTTPConfigs     *HTTPConfigs
	}
)

var dotEnvConfig = dotenv.Configure

func New() *ConfigBuilderImpl {
	return &ConfigBuilderImpl{}
}

func (b *ConfigBuilderImpl) Build() (*Configs, error) {
	appConfigs, err := b.getAppConfigs()
	if err != nil {
		return nil, err
	}

	sqlDatabaseConfigs, err := b.getSqlDatabaseConfigs()
	if err != nil {
		return nil, err
	}

	rabbitMQConfigs, err := b.getRabbitMQConfigs()
	if err != nil {
		return nil, err
	}

	otelConfigs, err := b.getOtelConfigs()
	if err != nil {
		return nil, err
	}

	httpServerConfigs, err := b.getHTTPServerConfigs()
	if err != nil {
		return nil, err
	}

	return &Configs{
		AppConfigs:      appConfigs,
		SqlConfigs:      sqlDatabaseConfigs,
		RabbitMqConfigs: rabbitMQConfigs,
		OtelConfigs:     otelConfigs,
		HTTPConfigs:     httpServerConfigs,
	}, nil
}

func (b *ConfigBuilderImpl) getAppConfigs() (*AppConfigs, error) {
	configs := AppConfigs{}
	configs.GoEnv = NewEnvironment(os.Getenv(GO_ENV_KEY))

	if configs.GoEnv == UNKNOWN_ENV {
		return nil, errors.New("[ConfigBuilder::New] unknown env")
	}

	err := dotEnvConfig(".env." + configs.GoEnv.ToString())
	if err != nil {
		return nil, err
	}

	configs.LogLevel = NewLogLevel(os.Getenv(LOG_LEVEL_ENV_KEY))
	configs.AppName = b.appName()

	return &configs, nil
}

func (b *ConfigBuilderImpl) appName() string {
	name := os.Getenv(APP_NAME_ENV_KEY)

	if name == "" {
		return DEFAULT_APP_NAME
	}

	return name
}

func (b *ConfigBuilderImpl) logPath(appName string) string {
	relative := os.Getenv(LOG_PATH_ENV_KEY)

	projectPath, _ := os.Getwd()

	if relative == "" {
		return fmt.Sprintf("%s%s%s%s", projectPath, DEFAULT_LOG_PATH, appName, ".log")
	}

	if relative[:1] == "." {
		relative = relative[1:]
	}

	separator := ""
	if len(relative) >= 1 && relative[:1] != string(os.PathSeparator) {
		separator = string(os.PathSeparator)
	}

	return fmt.Sprintf("%s%s%s", projectPath, separator, relative)
}
