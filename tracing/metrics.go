package tracing

import (
	"context"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
)

type (
	Headers         map[string]string
	ExporterType    int32
	OTLPCompression string

	TracingBuilder interface {
		AddHeader(key, value string) TracingBuilder
		WithHeaders(headers Headers) TracingBuilder
		Type(t ExporterType) TracingBuilder
		Endpoint(s string) TracingBuilder
		Build() (shutdown func(context.Context) error, err error)
	}

	OTLPTracingBuilder interface {
		TracingBuilder
		WithApiKeyHeader() OTLPTracingBuilder
		WithTimeout(t time.Duration) OTLPTracingBuilder
		WithReconnection(t time.Duration) OTLPTracingBuilder
		WithCompression(c OTLPCompression) OTLPTracingBuilder
	}

	JaegerTracingBuilder interface {
		TracingBuilder
	}

	tracingBuilder struct {
		logger logging.Logger
		cfg    *env.Configs

		headers      Headers
		exporterType ExporterType
		endpoint     string
	}

	otlpTracingBuilder struct {
		tracingBuilder
		reconnectionPeriod time.Duration
		timeout            time.Duration
		compression        OTLPCompression
	}

	jaegerTracingBuilder struct {
		tracingBuilder
	}
)

const (
	UNKNOWN_EXPORTER       ExporterType = 0
	OTLP_TLS_GRPC_EXPORTER ExporterType = 1
	OTLP_GRPC_EXPORTER     ExporterType = 2
	OTLP_HTTPS_EXPORTER    ExporterType = 3
	JAEGER_EXPORTER        ExporterType = 4

	OTLP_GZIP_COMPRESSIONS OTLPCompression = "gzip"
)

func Message(msg string) string {
	return "[gokit::tracing] " + msg
}
