// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package errors provides custom error types and error creation utilities
// specific to the configs_builder module. These errors help in providing
// clear information about configuration issues.
package errors

import "fmt"

// ConfigsError is a custom error type for configuration-related errors
type ConfigsError struct {
	msg string
}

// NewConfigsError creates a new ConfigsError with the given message,
// prefixing it with "configs builder error - "
func NewConfigsError(msg string) error {
	return &ConfigsError{
		msg: fmt.Sprintf("configs builder error - %s", msg),
	}
}

// Error implements the error interface for ConfigsError
func (e *ConfigsError) Error() string {
	return e.msg
}

// NewErrRequiredConfig creates an error for a missing required configuration
// parameter, using the provided environment variable name in the message
func NewErrRequiredConfig(env string) error {
	return NewConfigsError(fmt.Sprintf("%s is required", env))
}

// Pre-defined errors
var (
	// ErrUnknownEnv indicates that the application environment could not be determined
	ErrUnknownEnv = NewConfigsError("unknown env")
)
