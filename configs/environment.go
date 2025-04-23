// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// Environment represents the application execution environment as a typed enum.
// It provides type safety when dealing with different deployment environments.
type Environment int8

const (
	// UnknownEnv represents an unspecified or invalid environment
	UnknownEnv Environment = 0
	// LocalEnv represents a local development environment
	LocalEnv Environment = 1
	// DevelopmentEnv represents a shared development environment
	DevelopmentEnv Environment = 2
	// StagingEnv represents a pre-production environment for testing
	StagingEnv Environment = 3
	// QaEnv represents a quality assurance testing environment
	QaEnv Environment = 4
	// ProductionEnv represents the live production environment
	ProductionEnv Environment = 5
)

var (
	// EnvironmentMapping maps Environment enum values to their string representations
	EnvironmentMapping = map[Environment]string{
		UnknownEnv:     "unknown",
		LocalEnv:       "local",
		DevelopmentEnv: "development",
		StagingEnv:     "staging",
		QaEnv:          "qa",
		ProductionEnv:  "production",
	}
)

// NewEnvironment converts a string environment name to the corresponding Environment enum value.
// It accepts various standard naming conventions like "dev"/"development", case-insensitive.
// Returns UnknownEnv if the string doesn't match any known environment.
func NewEnvironment(env string) Environment {
	switch env {
	case "local":
		fallthrough
	case "LOCAL":
		return LocalEnv
	case "development":
		fallthrough
	case "DEVELOPMENT":
		fallthrough
	case "dev":
		fallthrough
	case "DEV":
		return DevelopmentEnv
	case "production":
		fallthrough
	case "PRODUCTION":
		fallthrough
	case "prod":
		fallthrough
	case "PROD":
		return ProductionEnv
	case "staging":
		fallthrough
	case "STAGING":
		fallthrough
	case "stg":
		fallthrough
	case "STG":
		return StagingEnv
	case "qa":
		fallthrough
	case "QA":
		return QaEnv
	default:
		return UnknownEnv
	}
}

// ToString returns the canonical string representation of the Environment.
// This provides consistent environment naming when logging or presenting the environment.
func (e Environment) ToString() string {
	switch e {
	case LocalEnv:
		return "local"
	case DevelopmentEnv:
		return "development"
	case ProductionEnv:
		return "production"
	case StagingEnv:
		return "staging"
	case QaEnv:
		return "qa"
	default:
		return "unknown"
	}
}
