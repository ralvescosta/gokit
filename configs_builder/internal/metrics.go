package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadMetricsConfigs() (*configs.MetricsConfigs, error) {
	enabled := os.Getenv(keys.METRICS_ENABLED_ENV_KEY)

	metricsConfigs := configs.MetricsConfigs{}

	if enabled != "" || enabled == "true" {
		metricsConfigs.Enabled = true
	}

	metricsConfigs.OtlpEndpoint = os.Getenv(keys.METRICS_OTLP_ENDPOINT_ENV_KEY)
	metricsConfigs.OtlpAPIKey = os.Getenv(keys.METRICS_OTLP_API_KEY_ENV_KEY)

	kind := os.Getenv(keys.METRICS_KIND_ENV_KEY)
	if kind == "prometheus" {
		metricsConfigs.Kind = configs.Prometheus
	}

	if kind == "otlp" {
		metricsConfigs.Kind = configs.OTLP
	}

	return &metricsConfigs, nil
}
