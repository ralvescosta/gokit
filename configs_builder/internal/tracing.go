// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

// ReadTracingConfigs retrieves OpenTelemetry tracing configuration from environment variables.
// Determines whether tracing is enabled and loads endpoint information for trace collection.
func ReadTracingConfigs() (*configs.TracingConfigs, error) {
	enabled := os.Getenv(keys.TracingEnabledEnvKey)

	configs := configs.TracingConfigs{}

	// Check if tracing is enabled
	if enabled != "" || enabled == "true" {
		configs.Enabled = true
	}

	// Load OTLP (OpenTelemetry Protocol) endpoint and API key
	configs.OtlpEndpoint = os.Getenv(keys.TracingOtlpEndpointEnvKey)
	configs.OtlpAPIKey = os.Getenv(keys.TracingOtlpAPIKeyEnvKey)

	return &configs, nil
}
