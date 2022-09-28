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
		Build(context.Context) (shutdown func(context.Context) error, err error)
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
		logger logging.ILogger
		cfg    *env.Config

		appName      string
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
