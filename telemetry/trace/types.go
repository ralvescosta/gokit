package trace

import (
	"context"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
)

type (
	Headers          map[string]string
	OTLPExporterType int32
	OTLPCompression  string

	TraceBuilder interface {
		WithApiKeyHeader() TraceBuilder
		AddHeader(key, value string) TraceBuilder
		WithHeaders(headers Headers) TraceBuilder
		Type(t OTLPExporterType) TraceBuilder
		Endpoint(s string) TraceBuilder
		WithTimeout(t time.Duration) TraceBuilder
		WithReconnection(t time.Duration) TraceBuilder
		WithCompression(c OTLPCompression) TraceBuilder
		Build(context.Context) (shutdown func(context.Context) error, err error)
	}

	traceBuilder struct {
		logger logging.ILogger
		cfg    *env.Configs

		appName            string
		headers            Headers
		exporterType       OTLPExporterType
		endpoint           string
		reconnectionPeriod time.Duration
		timeout            time.Duration
		compression        OTLPCompression
	}
)
