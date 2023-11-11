package configs

type MetricsKind string

const (
	OTLP       = MetricsKind("OTLP")
	Prometheus = MetricsKind("Prometheus")
)

type MetricsConfigs struct {
	Enabled bool

	Kind MetricsKind

	OtlpEndpoint string
	OtlpApiKey   string
}
