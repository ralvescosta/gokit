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

	MESSAGING_ENGINES_ENV_KEY = "MESSAGING_ENGINE_ENV_KEY"
	RABBIT_HOST_ENV_KEY       = "RABBIT_HOST_ENV_KEY"
	RABBIT_PORT_ENV_KEY       = "RABBIT_PORT_ENV_KEY"
	RABBIT_USER_ENV_KEY       = "RABBIT_USER_ENV_KEY"
	RABBIT_PASSWORD_ENV_KEY   = "RABBIT_PASSWORD_ENV_KEY"
	RABBIT_VHOST_ENV_KEY      = "RABBIT_VHOST_ENV_KEY"
	KAFKA_HOST_ENV_KEY        = "KAFKA_HOST_ENV_KEY"
	KAFKA_PORT_ENV_KEY        = "KAFKA_PORT_ENV_KEY"
	KAFKA_USER_ENV_KEY        = "KAFKA_USER_ENV_KEY"
	KAFKA_PASSWORD_ENV_KEY    = "KAFKA_PASSWORD_ENV_KEY"
	RABBITMQ_ENGINE           = "RabbitMQ"
	KAFKA_ENGINE              = "Kafka"

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

	IS_TRACING_ENABLED_ENV_KEY    = "TRACING_ENABLED"
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
		Database() ConfigBuilder
		Messaging() ConfigBuilder
		Tracing() ConfigBuilder
		HTTPServer() ConfigBuilder
		Build() (*Config, error)
	}

	Config struct {
		Err error

		GO_ENV Environment

		LOG_LEVEL LogLevel
		LOG_PATH  string

		APP_NAME string

		SQL_DB_HOST            string
		SQL_DB_PORT            string
		SQL_DB_USER            string
		SQL_DB_PASSWORD        string
		SQL_DB_NAME            string
		SQL_DB_SECONDS_TO_PING int

		MESSAGING_ENGINES map[string]bool
		RABBIT_HOST       string
		RABBIT_PORT       string
		RABBIT_USER       string
		RABBIT_PASSWORD   string
		RABBIT_VHOST      string
		KAFKA_HOST        string
		KAFKA_PORT        string
		KAFKA_USER        string
		KAFKA_PASSWORD    string

		IS_TRACING_ENABLED        bool
		OTLP_ENDPOINT             string
		OTLP_API_KEY              string
		JAEGER_SERVICE_NAME       string
		JAEGER_AGENT_HOST         string
		JAEGER_SAMPLER_TYPE       string
		JAEGER_SAMPLER_PARAM      int
		JAEGER_REPORTER_LOG_SPANS bool
		JAEGER_RPC_METRICS        bool

		HTTP_PORT string
		HTTP_HOST string
		HTTP_ADDR string
	}
)

var dotEnvConfig = dotenv.Configure

func New() ConfigBuilder {
	c := &Config{}

	c.GO_ENV = NewEnvironment(os.Getenv(GO_ENV_KEY))

	if c.GO_ENV == UNKNOWN_ENV {
		c.Err = errors.New("[ConfigBuilder::New] unknown env")
		return c
	}

	err := dotEnvConfig(".env." + EnvironmentMapping[c.GO_ENV])
	if err != nil {
		c.Err = err
		return c
	}

	return c
}

func (c *Config) Build() (*Config, error) {
	if c.Err != nil {
		return c, c.Err
	}

	c.LOG_LEVEL = NewLogLevel(os.Getenv(LOG_LEVEL_ENV_KEY))
	c.APP_NAME = NewAppName()
	c.LOG_PATH = NewLogPath(c.APP_NAME)

	return c, nil
}

func NewAppName() string {
	name := os.Getenv(APP_NAME_ENV_KEY)

	if name == "" {
		return DEFAULT_APP_NAME
	}

	return name
}

func NewLogPath(appName string) string {
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
