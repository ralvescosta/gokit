package configs

type LogLevel int8

const (
	DEBUG LogLevel = 0
	INFO  LogLevel = 1
	WARN  LogLevel = 2
	ERROR LogLevel = 3
	PANIC LogLevel = 4
)

func NewLogLevel(env string) LogLevel {
	switch env {
	case "debug":
		fallthrough
	case "DEBUG":
		return DEBUG
	case "warn":
		fallthrough
	case "WARN":
		return INFO
	case "error":
		fallthrough
	case "ERROR":
		return ERROR
	case "panic":
		fallthrough
	case "PANIC":
		return PANIC
	default:
		return INFO
	}
}
