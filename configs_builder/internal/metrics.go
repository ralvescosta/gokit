package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadMetricsConfigs() (*configs.MetricsConfigs, error) {
	enabled := os.Getenv(keys.MetricsEnabledEnvKey)

	metricsConfigs := configs.MetricsConfigs{}

	if enabled != "" || enabled == "true" {
		metricsConfigs.Enabled = true
	}

	metricsConfigs.OtlpEndpoint = os.Getenv(keys.MetricsOtlpEndpointEnvKey)
	metricsConfigs.OtlpAPIKey = os.Getenv(keys.MetricsOtlpAPIKeyEnvKey)

	kind := os.Getenv(keys.MetricsKindEnvKey)
	if kind == "prometheus" {
		metricsConfigs.Kind = configs.Prometheus
	}

	if kind == "otlp" {
		metricsConfigs.Kind = configs.OTLP
	}

	return &metricsConfigs, nil
}
