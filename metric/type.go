package metric

import (
	"context"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
)

type (
	Headers            map[string]string
	MetricExporterType int32
	OTLPCompression    string

	MetricBuilder interface {
		WithApiKeyHeader() MetricBuilder
		AddHeader(key, value string) MetricBuilder
		WithHeaders(headers Headers) MetricBuilder
		Endpoint(s string) MetricBuilder
		WithTimeout(t time.Duration) MetricBuilder
		WithReconnection(t time.Duration) MetricBuilder
		WithCompression(c OTLPCompression) MetricBuilder
		Build(context.Context) (shutdown func(context.Context) error, err error)
	}

	metricBuilder struct {
		logger logging.Logger
		cfg    *env.Config

		appName            string
		headers            Headers
		endpoint           string
		reconnectionPeriod time.Duration
		timeout            time.Duration
		compression        OTLPCompression
	}

	otlpMetricBuilder struct {
		metricBuilder
	}

	prometheusMetricBuilder struct {
		metricBuilder
	}
)
