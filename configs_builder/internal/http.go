// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package internal

import (
	"fmt"
	"os"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

// ReadHTTPConfigs retrieves and validates HTTP server configuration from environment variables.
// Returns an error if required configuration values are missing.
func ReadHTTPConfigs() (*configs.HTTPConfigs, error) {
	httpConfigs := configs.HTTPConfigs{}

	// Get and validate HTTP port
	httpConfigs.Port = os.Getenv(keys.HTTPPortEnvKey)
	if httpConfigs.Port == "" {
		return nil, errors.NewErrRequiredConfig(keys.HTTPPortEnvKey)
	}

	// Get and validate HTTP host
	httpConfigs.Host = os.Getenv(keys.HTTPHostEnvKey)
	if httpConfigs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.HTTPHostEnvKey)
	}

	// Construct full address string from host and port
	httpConfigs.Addr = fmt.Sprintf("%s:%s", httpConfigs.Host, httpConfigs.Port)

	// Check if profiling is enabled
	profiling := os.Getenv(keys.HTTPEnableProfilingEnvKey)
	if profiling == "true" {
		httpConfigs.EnableProfiling = true
	}

	return &httpConfigs, nil
}
