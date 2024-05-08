package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadTracingConfigs() (*configs.TracingConfigs, error) {
	enabled := os.Getenv(keys.TRACING_ENABLED_ENV_KEY)

	configs := configs.TracingConfigs{}

	if enabled != "" || enabled == "true" {
		configs.Enabled = true
	}

	configs.OtlpEndpoint = os.Getenv(keys.TRACING_OTLP_ENDPOINT_ENV_KEY)
	configs.OtlpAPIKey = os.Getenv(keys.TRACING_OTLP_API_KEY_ENV_KEY)

	return &configs, nil
}
