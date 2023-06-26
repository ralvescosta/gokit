package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadOtelConfigs() (*configs.OtelConfigs, error) {
	tracingEnabled := os.Getenv(keys.TRACING_ENABLED_ENV_KEY)
	metricsEnabled := os.Getenv(keys.METRICS_ENABLED_ENV_KEY)

	if tracingEnabled == "" && metricsEnabled == "" {
		return nil, errors.NewErrRequiredConfig(keys.TRACING_ENABLED_ENV_KEY)
	}

	otelConfigs := configs.OtelConfigs{}

	if tracingEnabled == "true" {
		otelConfigs.TracingEnabled = true
	}

	if metricsEnabled == "true" {
		otelConfigs.MetricsEnabled = true
	}

	otelConfigs.OtlpEndpoint = os.Getenv(keys.OTLP_ENDPOINT_ENV_KEY)
	otelConfigs.OtlpApiKey = os.Getenv(keys.OTLP_API_KEY_ENV_KEY)
	otelConfigs.JaegerServiceName = os.Getenv(keys.JAEGER_SERVICE_NAME_KEY)
	otelConfigs.JaegerAgentHost = os.Getenv(keys.JAEGER_AGENT_HOST_KEY)
	otelConfigs.JaegerSampleType = os.Getenv(keys.JAEGER_SAMPLER_TYPE_KEY)
	if samplerParam := os.Getenv(keys.JAEGER_SAMPLER_PARAM_KEY); samplerParam != "" {
		otelConfigs.JaegerSampleParam, _ = strconv.Atoi(samplerParam)
	}

	if reportLogSpans := os.Getenv(keys.JAEGER_REPORTER_LOG_SPANS_KEY); reportLogSpans != "" {
		otelConfigs.JaegerReporterLogSpans = reportLogSpans == "true"
	}

	if rpcMetrics := os.Getenv(keys.JAEGER_RPC_METRICS_KEY); rpcMetrics != "" {
		otelConfigs.JaegerRpcMetrics = rpcMetrics == "true"
	}

	return &otelConfigs, nil
}
