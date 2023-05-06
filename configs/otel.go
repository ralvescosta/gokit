package configs

type OtelConfigs struct {
	TracingEnabled bool
	MetricsEnabled bool

	OtlpEndpoint string
	OtlpApiKey   string

	JaegerServiceName      string
	JaegerAgentHost        string
	JaegerSampleType       string
	JaegerSampleParam      int
	JaegerReporterLogSpans bool
	JaegerRpcMetrics       bool
}
