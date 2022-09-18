package tracing

const (
	UNKNOWN_EXPORTER       ExporterType = 0
	OTLP_TLS_GRPC_EXPORTER ExporterType = 1
	OTLP_GRPC_EXPORTER     ExporterType = 2
	OTLP_HTTPS_EXPORTER    ExporterType = 3
	JAEGER_EXPORTER        ExporterType = 4

	OTLP_GZIP_COMPRESSIONS OTLPCompression = "gzip"
)

func LogMessage(msg string) string {
	return "[gokit::tracing] " + msg
}
