package env

type Environment int8
type LogLevel int8

const (
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
)

type (
	Env struct {
		GO_ENV    Environment
		LOG_LEVEL LogLevel
	}
)
