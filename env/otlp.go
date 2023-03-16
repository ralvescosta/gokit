package env

import (
	"os"
	"strconv"
)

type (
	OtelConfigs struct {
		TracingEnabled         bool
		MetricsEnabled         bool
		OtlpEndpoint           string
		OtlpApiKey             string
		JaegerServiceName      string
		JaegerAgentHost        string
		JaegerSampleType       string
		JaegerSampleParam      int
		JaegerReporterLogSpans bool
		JaegerRpcMetrics       bool
	}
)

func (b *ConfigsBuilderImpl) Otel() ConfigsBuilder {
	b.otel = true
	return b
}

func (b *ConfigsBuilderImpl) getOtelConfigs() (*OtelConfigs, error) {
	if !b.otel {
		return nil, nil
	}

	tracingEnabled := os.Getenv(TRACING_ENABLED_ENV_KEY)
	metricsEnabled := os.Getenv(METRICS_ENABLED_ENV_KEY)

	if tracingEnabled == "" || metricsEnabled == "" {
		return nil, NewErrRequiredConfig(TRACING_ENABLED_ENV_KEY)
	}

	configs := OtelConfigs{}

	if tracingEnabled == "true" {
		configs.TracingEnabled = true
	}

	if metricsEnabled == "true" {
		configs.MetricsEnabled = true
	}

	configs.OtlpEndpoint = os.Getenv(OTLP_ENDPOINT_ENV_KEY)
	configs.OtlpApiKey = os.Getenv(OTLP_API_KEY_ENV_KEY)
	configs.JaegerServiceName = os.Getenv(JAEGER_SERVICE_NAME_KEY)
	configs.JaegerAgentHost = os.Getenv(JAEGER_AGENT_HOST_KEY)
	configs.JaegerSampleType = os.Getenv(JAEGER_SAMPLER_TYPE_KEY)
	if samplerParam := os.Getenv(JAEGER_SAMPLER_PARAM_KEY); samplerParam != "" {
		configs.JaegerSampleParam, _ = strconv.Atoi(samplerParam)
	}

	if reportLogSpans := os.Getenv(JAEGER_REPORTER_LOG_SPANS_KEY); reportLogSpans != "" {
		configs.JaegerReporterLogSpans = reportLogSpans == "true"
	}

	if rpcMetrics := os.Getenv(JAEGER_RPC_METRICS_KEY); rpcMetrics != "" {
		configs.JaegerRpcMetrics = rpcMetrics == "true"
	}

	return &configs, nil
}
