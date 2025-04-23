// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package metrics provides tools for collecting, exporting, and monitoring metrics
// with support for multiple exporters like OpenTelemetry and Prometheus.
package metrics

// MetricExporterType represents the type of metrics exporter used in the application.
const (
	// UnknownExporter represents an unspecified metrics exporter type.
	UnknownExporter MetricExporterType = 0
	// OtlpGrpcTLSExporter represents the OpenTelemetry gRPC exporter with TLS enabled.
	OtlpGrpcTLSExporter MetricExporterType = 1
	// OtlpGrpcExporter represents the standard OpenTelemetry gRPC exporter.
	OtlpGrpcExporter MetricExporterType = 2
	// PrometheusExporter represents the Prometheus metrics exporter.
	PrometheusExporter MetricExporterType = 3

	// OtlpGzipCompressions represents gzip compression for OTLP exporters.
	OtlpGzipCompressions OTLPCompression = "gzip"
)

// Message creates a standardized prefix for metrics-related log messages.
// This helps to identify metrics-related logs in the application output.
//
// Parameters:
//   - msg: The message text to be prefixed.
//
// Returns:
//   - A formatted string with the metrics package prefix.
func Message(msg string) string {
	return "[gokit:metrics] " + msg
}
