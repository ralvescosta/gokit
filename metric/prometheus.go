package metric

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
)

func NewPrometheus(cfg *env.Config, logger logging.ILogger) MetricBuilder {
	return &prometheusMetricBuilder{
		metricBuilder: metricBuilder{
			logger:             logger,
			cfg:                cfg,
			appName:            cfg.APP_NAME,
			endpoint:           cfg.OTLP_ENDPOINT,
			reconnectionPeriod: 2 * time.Second,
			timeout:            30 * time.Second,
			compression:        OTLP_GZIP_COMPRESSIONS,
			headers:            Headers{},
		},
	}
}

func (b *prometheusMetricBuilder) WithApiKeyHeader() MetricBuilder {
	b.headers["api-key"] = b.cfg.OTLP_API_KEY
	return b
}

func (b *prometheusMetricBuilder) AddHeader(key, value string) MetricBuilder {
	b.headers[key] = value
	return b
}

func (b *prometheusMetricBuilder) WithHeaders(headers Headers) MetricBuilder {
	b.headers = headers
	return b
}

func (b *prometheusMetricBuilder) Endpoint(s string) MetricBuilder {
	b.endpoint = s
	return b
}

func (b *prometheusMetricBuilder) WithTimeout(t time.Duration) MetricBuilder {
	b.timeout = t
	return b
}

func (b *prometheusMetricBuilder) WithReconnection(t time.Duration) MetricBuilder {
	b.reconnectionPeriod = t
	return b
}

func (b *prometheusMetricBuilder) WithCompression(c OTLPCompression) MetricBuilder {
	b.compression = c
	return b
}

func (b *prometheusMetricBuilder) Build(ctx context.Context) (shutdown func(context.Context) error, err error) {
	return b.prometheusExporter(ctx)
}

func (b *prometheusMetricBuilder) prometheusExporter(ctx context.Context) (shutdown func(context.Context) error, err error) {
	b.logger.Debug(LogMessage("prometheus metric exporter"))

	b.logger.Debug(LogMessage("creating prometheus resource..."))
	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("library.language", "go"),
			attribute.String("service.name", b.appName),
			attribute.String("environment", b.cfg.GO_ENV.ToString()),
			attribute.Int64("ID", int64(os.Getegid())),
		),
	)

	if err != nil {
		b.logger.Error(LogMessage("could not set resources"), logging.ErrorField(err))
		return nil, err
	}
	b.logger.Debug(LogMessage("prometheus resource created"))

	b.logger.Debug(LogMessage("configuring prometheus provider..."))
	pConfig := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}
	metricProvider := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(pConfig.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(resources),
	)
	b.logger.Debug(LogMessage("prometheus provider configured"))

	b.logger.Debug(LogMessage("starting prometheus provider..."))
	exporter, err := prometheus.New(pConfig, metricProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize prometheus exporter: %w", err)
	}
	global.SetMeterProvider(exporter.MeterProvider())
	b.logger.Debug(LogMessage("prometheus provider started"))

	b.logger.Debug(LogMessage("prometheus metric exporter configured"))
	return metricProvider.Stop, nil
}
