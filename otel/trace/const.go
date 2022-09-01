package trace

const (
	UNKNOWN_EXPORTER  OTLPExporterType = 0
	TLS_GRPC_EXPORTER OTLPExporterType = 1
	GRPC_EXPORTER     OTLPExporterType = 2
	HTTPS_EXPORTER    OTLPExporterType = 3
	HTTP_EXPORTER     OTLPExporterType = 4

	OTLP_GZIP_COMPRESSIONS OTLPCompression = "gzip"
)

func LogMessage(msg string) string {
	return "[gokit::otel::trace] " + msg
}
