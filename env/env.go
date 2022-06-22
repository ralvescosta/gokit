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
	IEnv interface {
		Load() error
	}

	Env struct {
		GO_ENV    Environment
		LOG_LEVEL LogLevel
		LOG_PATH  string
		APP_NAME  string
	}
)

var dotEnvConfig = dotenv.Configure

func (e *Env) Load() error {
	if e.GO_ENV == UNKNOWN_ENV {
		return errors.New("unknown env")
	}
	err := dotEnvConfig(EnvironmentMapping[e.GO_ENV])
	if err != nil {
		return err
	}

	return nil
}

func NewAppEnvironment() IEnv {
	goEnv := NewEnvironment(os.Getenv(GO_ENV_KEY))
	logLevel := NewLogLevel(os.Getenv(LOG_LEVEL_ENV_KEY))
	appName := NewAppName()
	logPath := NewLogPath(appName)

	return &Env{
		GO_ENV:    goEnv,
		LOG_LEVEL: logLevel,
		APP_NAME:  appName,
		LOG_PATH:  logPath,
	}
}

type Environment int8

func NewEnvironment(env string) Environment {
	switch env {
	case "development":
		fallthrough
	case "DEVELOPMENT":
		fallthrough
	case "dev":
		fallthrough
	case "DEV":
		return DEVELOPMENT_ENV
	case "production":
		fallthrough
	case "PRODUCTION":
		fallthrough
	case "prod":
		fallthrough
	case "PROD":
		return PRODUCTION_ENV
	case "staging":
		fallthrough
	case "STAGING":
		fallthrough
	case "stg":
		fallthrough
	case "STG":
		return STAGING_ENV
	case "qa":
		fallthrough
	case "QA":
		return QA_ENV
	default:
		return UNKNOWN_ENV
	}
}

type LogLevel int8

func NewLogLevel(env string) LogLevel {
	switch env {
	case "debug":
		fallthrough
	case "DEBUG":
		return DEBUG_L
	case "warn":
		fallthrough
	case "WARN":
		return WARN_L
	case "error":
		fallthrough
	case "ERROR":
		return ERROR_L
	case "panic":
		fallthrough
	case "PANIC":
		return PANIC_L
	default:
		return INFO_L
	}
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
