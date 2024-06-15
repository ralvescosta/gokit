package metrics

const (
	UnknownExporter     MetricExporterType = 0
	OtlpGrpcTLSExporter MetricExporterType = 1
	OtlpGrpcExporter    MetricExporterType = 2
	PrometheusExporter  MetricExporterType = 3

	OtlpGzipCompressions OTLPCompression = "gzip"
)

func Message(msg string) string {
	return "[gokit:metrics] " + msg
}
