package trace

const (
	UNKNOWN_EXPORTER OTLPExporterType = 0
	GRPC_EXPORTER    OTLPExporterType = 1
	HTTPS_EXPORTER   OTLPExporterType = 2

	OTLP_GZIP_COMPRESSIONS OTLPCompression = "gzip"
)
