// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package tracing provides distributed tracing capabilities using OpenTelemetry
// to help monitor and troubleshoot microservices-based applications.
package tracing

import (
	"context"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
)

type (
	// Headers represents a map of string key-value pairs used for propagating trace context
	Headers map[string]string

	// ExporterType defines the type of the tracing exporter to be used
	ExporterType int32

	// OTLPCompression defines the compression algorithm used for OTLP exporter
	OTLPCompression string

	// TracingBuilder provides an interface for configuring and building a tracing provider
	// using the builder pattern. It allows setting headers, exporter type, and endpoint.
	TracingBuilder interface {
		// AddHeader adds a single header key-value pair to the exporter configuration
		AddHeader(key, value string) TracingBuilder

		// WithHeaders sets all headers at once for the exporter configuration
		WithHeaders(headers Headers) TracingBuilder

		// Type sets the exporter type to be used
		Type(t ExporterType) TracingBuilder

		// Endpoint sets the endpoint URL for the exporter
		Endpoint(s string) TracingBuilder

		// Build creates and configures the tracing provider based on the builder settings
		// Returns a shutdown function to cleanly close the exporter and any error encountered
		Build() (shutdown func(context.Context) error, err error)
	}

	// tracingBuilder is the internal implementation of the TracingBuilder interface
	tracingBuilder struct {
		logger logging.Logger
		cfg    *configs.Configs

		headers      Headers
		exporterType ExporterType
		endpoint     string
	}
)

const (
	// UNKNOWN_EXPORTER represents an undefined exporter type
	UNKNOWN_EXPORTER ExporterType = 0

	// OTLP_TLS_GRPC_EXPORTER represents the secure OTLP gRPC exporter type (with TLS)
	OTLP_TLS_GRPC_EXPORTER ExporterType = 1

	// OTLP_GRPC_EXPORTER represents the insecure OTLP gRPC exporter type (without TLS)
	OTLP_GRPC_EXPORTER ExporterType = 2

	// OTLP_HTTPS_EXPORTER represents the OTLP HTTPS exporter type
	OTLP_HTTPS_EXPORTER ExporterType = 3

	// JAEGER_EXPORTER represents the Jaeger exporter type
	JAEGER_EXPORTER ExporterType = 4

	// OTLP_GZIP_COMPRESSIONS represents gzip compression for OTLP exporters
	OTLP_GZIP_COMPRESSIONS OTLPCompression = "gzip"
)

// Message prefixes a string with the package identifier for standardized logging
func Message(msg string) string {
	return "[gokit::tracing] " + msg
}
