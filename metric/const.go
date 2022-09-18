package metric

const (
	UNKNOWN_EXPORTER       MetricExporterType = 0
	OTLP_GRPC_TLS_EXPORTER MetricExporterType = 1
	OTLP_GRPC_EXPORTER     MetricExporterType = 2
	PROMETHEUS_EXPORTER    MetricExporterType = 3

	OTLP_GZIP_COMPRESSIONS OTLPCompression = "gzip"
)

func LogMessage(msg string) string {
	return "[gokit::metric] " + msg
}
