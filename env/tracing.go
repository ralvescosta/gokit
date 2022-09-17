package env

import (
	"fmt"
	"os"
	"strconv"
)

const (
	RequiredTelemetryErrorMessage = "[ConfigBuilder::Telemetry] %s is required"
)

func (c *Config) Tracing() ConfigBuilder {
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
	c.OTLP_API_KEY = os.Getenv(OTLP_API_KEY_ENV_KEY)
	c.JAEGER_SERVICE_NAME = os.Getenv(JAEGER_SERVICE_NAME_KEY)
	c.JAEGER_AGENT_HOST = os.Getenv(JAEGER_AGENT_HOST_KEY)
	c.JAEGER_SAMPLER_TYPE = os.Getenv(JAEGER_SAMPLER_TYPE_KEY)
	if samplerParam := os.Getenv(JAEGER_SAMPLER_PARAM_KEY); samplerParam != "" {
		c.JAEGER_SAMPLER_PARAM, _ = strconv.Atoi(samplerParam)
	}
	if reportLogSpans := os.Getenv(JAEGER_REPORTER_LOG_SPANS_KEY); reportLogSpans != "" {
		c.JAEGER_REPORTER_LOG_SPANS = reportLogSpans == "true"
	}
	if rpcMetrics := os.Getenv(JAEGER_RPC_METRICS_KEY); rpcMetrics != "" {
		c.JAEGER_RPC_METRICS = rpcMetrics == "true"
	}

	return c
}
