// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

// ReadSQLDatabaseConfigs retrieves SQL database configuration from environment variables.
// Validates that all required database connection parameters are provided and returns
// an error if any required configuration is missing.
func ReadSQLDatabaseConfigs() (*configs.SQLConfigs, error) {
	sqlConfigs := configs.SQLConfigs{}

	// Get and validate database host
	sqlConfigs.Host = os.Getenv(keys.SQLDbHostEnvKey)
	if sqlConfigs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQLDbHostEnvKey)
	}

	// Get and validate database port
	sqlConfigs.Port = os.Getenv(keys.SQLDbPortEnvKey)
	if sqlConfigs.Port == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQLDbPortEnvKey)
	}

	// Get and validate database user
	sqlConfigs.User = os.Getenv(keys.SQLDbUserEnvKey)
	if sqlConfigs.User == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQLDbUserEnvKey)
	}

	// Get and validate database password
	sqlConfigs.Password = os.Getenv(keys.SQLDbPasswordEnvKey)
	if sqlConfigs.Password == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQLDbPasswordEnvKey)
	}

	// Get and validate database name
	sqlConfigs.DbName = os.Getenv(keys.SQLDbNameEnvKey)
	if sqlConfigs.DbName == "" {
		return nil, errors.NewErrRequiredConfig(keys.SQLDbNameEnvKey)
	}

	// Parse and set the database ping interval
	p, err := strconv.Atoi(os.Getenv(keys.SQLDbSecondsToPingEnvKey))
	if err != nil {
		return nil, err
	}

	sqlConfigs.SecondsToPing = p

	return &sqlConfigs, nil
}
