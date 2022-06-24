package env

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
