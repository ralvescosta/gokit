// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

// ReadMetricsConfigs retrieves metrics collection configuration from environment variables.
// Determines whether metrics collection is enabled and loads endpoint information for metrics export.
func ReadMetricsConfigs() (*configs.MetricsConfigs, error) {
	enabled := os.Getenv(keys.MetricsEnabledEnvKey)

	metricsConfigs := configs.MetricsConfigs{}

	// Check if metrics collection is enabled
	if enabled != "" || enabled == "true" {
		metricsConfigs.Enabled = true
	}

	// Load OTLP (OpenTelemetry Protocol) endpoint and API key for metrics
	metricsConfigs.OtlpEndpoint = os.Getenv(keys.MetricsOtlpEndpointEnvKey)
	metricsConfigs.OtlpAPIKey = os.Getenv(keys.MetricsOtlpAPIKeyEnvKey)

	// Determine which metrics system to use based on environment configuration
	kind := os.Getenv(keys.MetricsKindEnvKey)
	if kind == "prometheus" {
		metricsConfigs.Kind = configs.Prometheus
	}

	if kind == "otlp" {
		metricsConfigs.Kind = configs.OTLP
	}

	return &metricsConfigs, nil
}
