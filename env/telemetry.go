package env

import (
	"fmt"
	"os"
)

const (
	RequiredTelemetryErrorMessage = "[ConfigBuilder::Telemetry] %s is required"
)

func (c *Configs) Tracing() IConfigs {
	if c.Err != nil {
		return c
	}

	tracingEnabled := os.Getenv(IS_TRACING_ENABLED_ENV_KEY)
	if tracingEnabled == "" {
		c.Err = fmt.Errorf(RequiredTelemetryErrorMessage, IS_TRACING_ENABLED_ENV_KEY)
		return c
	}

	if tracingEnabled == "true" {
		c.IS_TRACING_ENABLED = true
	}

	c.OTLP_ENDPOINT = os.Getenv(OTLP_ENDPOINT_ENV_KEY)
	if c.OTLP_ENDPOINT == "" {
		c.Err = fmt.Errorf(RequiredTelemetryErrorMessage, OTLP_ENDPOINT_ENV_KEY)
		return c
	}

	c.OTLP_API_KEY = os.Getenv(OTLP_API_KEY_ENV_KEY)

	return c
}
