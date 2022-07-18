package env

import (
	"errors"
	"fmt"
	"os"
)

const (
	RequiredTelemetryErrorMessage = "[ConfigBuilder::Telemetry] %s is required"
)

func (c *Configs) Tracing() IConfigs {
	tracingEnabled := os.Getenv(IS_TRACING_ENABLED_ENV_KEY)
	if tracingEnabled == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredTelemetryErrorMessage, IS_TRACING_ENABLED_ENV_KEY))
		return c
	}

	if tracingEnabled == "true" {
		c.IS_TRACING_ENABLED_ENV_KEY = true
		return c
	}

	c.IS_TRACING_ENABLED_ENV_KEY = false

	c.OTLP_ENDPOINT = os.Getenv(OTLP_ENDPOINT_ENV_KEY)
	if c.OTLP_ENDPOINT == "" {
		c.Err = errors.New(fmt.Sprintf(RequiredTelemetryErrorMessage, OTLP_ENDPOINT_ENV_KEY))
		return c
	}

	return c
}
