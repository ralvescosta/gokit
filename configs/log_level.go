// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// LogLevel defines the severity levels for application logging.
// It is used to filter log messages based on their importance.
type LogLevel int8

const (
	// DEBUG represents detailed information, typically valuable only for diagnosing problems
	DEBUG LogLevel = 0
	// INFO represents general information about system operation
	INFO LogLevel = 1
	// WARN represents potentially harmful situations that might need attention
	WARN LogLevel = 2
	// ERROR represents error events that might still allow the application to continue running
	ERROR LogLevel = 3
	// PANIC represents severe error events that will likely lead the application to abort
	PANIC LogLevel = 4
)

// NewLogLevel converts a string log level name to the corresponding LogLevel enum value.
// It accepts case-insensitive level names and defaults to INFO if the input doesn't match any known level.
func NewLogLevel(env string) LogLevel {
	switch env {
	case "debug":
		fallthrough
	case "DEBUG":
		return DEBUG
	case "warn":
		fallthrough
	case "WARN":
		return WARN
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
