package metric

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
)

type (
	Headers            map[string]string
	MetricExporterType int32
	OTLPCompression    string

	BasicMetricBuilder interface {
		Build() (shutdown func(context.Context) error, err error)
	}

	basicMetricBuilder struct {
		logger  logging.Logger
		cfg     *env.Config
		appName string
	}

	OTLPMetricBuilder interface {
		BasicMetricBuilder
		WithApiKeyHeader() OTLPMetricBuilder
		AddHeader(key, value string) OTLPMetricBuilder
		WithHeaders(headers Headers) OTLPMetricBuilder
		Endpoint(s string) OTLPMetricBuilder
		WithTimeout(t time.Duration) OTLPMetricBuilder
		WithReconnection(t time.Duration) OTLPMetricBuilder
		WithCompression(c OTLPCompression) OTLPMetricBuilder
	}

	otlpMetricBuilder struct {
		basicMetricBuilder

		headers            Headers
		endpoint           string
		reconnectionPeriod time.Duration
		timeout            time.Duration
		compression        OTLPCompression
	}

	PrometheusMetricBuilder interface {
		BasicMetricBuilder
		PromCollector() prometheus.Collector
	}

	prometheusMetricBuilder struct {
		basicMetricBuilder
		collector prometheus.Collector
	}
)
