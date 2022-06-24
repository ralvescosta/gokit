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

	SQL_DB_HOST_ENV_KEY     = "SQL_DB_HOST"
	SQL_DB_PORT_ENV_KEY     = "SQL_DB_PORT"
	SQL_DB_USER_ENV_KEY     = "SQL_DB_USER"
	SQL_DB_PASSWORD_ENV_KEY = "SQL_DB_PASSWORD"
	SQL_DB_NAME_ENV_KEY     = "SQL_DB_NAME"
	SQL_DB_SECONDS_TO_PING  = "SQL_DB_SECONDS_TO_PING"

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
	DEVELOPMENT_ENV Environment = 1
	STAGING_ENV     Environment = 2
	QA_ENV          Environment = 3
	PRODUCTION_ENV  Environment = 4

	DEBUG_L LogLevel = 0
	INFO_L  LogLevel = 1
	WARN_L  LogLevel = 2
	ERROR_L LogLevel = 3
	PANIC_L LogLevel = 4

	DEFAULT_APP_NAME = "app"
	DEFAULT_LOG_PATH = "/logs/"
)

var (
	EnvironmentMapping = map[Environment]string{
		UNKNOWN_ENV:     "unknown",
		DEVELOPMENT_ENV: "development",
		STAGING_ENV:     "staging",
		QA_ENV:          "qa",
		PRODUCTION_ENV:  "production",
	}
)

type (
	IConfigs interface {
		Database() IConfigs
		Messaging() IConfigs
		Build() (*Configs, error)
	}

	Configs struct {
		Err error

		GO_ENV Environment

		LOG_LEVEL LogLevel
		LOG_PATH  string

		APP_NAME string

		SQL_DB_HOST     string
		SQL_DB_PORT     string
		SQL_DB_USER     string
		SQL_DB_PASSWORD string
		SQL_DB_NAME     string

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
	}
)

var dotEnvConfig = dotenv.Configure

func New() IConfigs {
	c := &Configs{}

	c.GO_ENV = NewEnvironment(os.Getenv(GO_ENV_KEY))

	if c.GO_ENV == UNKNOWN_ENV {
		c.Err = errors.New("[ConfigBuilder::New] unknown env")
		return c
	}

	err := dotEnvConfig(EnvironmentMapping[c.GO_ENV])
	if err != nil {
		c.Err = err
		return c
	}

	return c
}

func (c *Configs) Build() (*Configs, error) {
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
