package metric

import (
	"context"
	"os"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel/attribute"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
)

func NewPrometheus(cfg *env.Config, logger logging.Logger) MetricBuilder {
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

//@TODO: Export the http handler to create prometheus scraping route
func (b *prometheusMetricBuilder) Build() (shutdown func(context.Context) error, err error) {

	b.logger.Debug(Message("prometheus metric exporter"))

	b.logger.Debug(Message("creating prometheus resource..."))
	ctx := context.Background()
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
		b.logger.Error(Message("could not set resources"), zap.Error(err))
		return nil, err
	}
	b.logger.Debug(Message("prometheus resource created"))

	b.logger.Debug(Message("starting prometheus provider..."))

	exporter := otelprom.New()
	provider := metric.NewMeterProvider(metric.WithReader(exporter), metric.WithResource(resources))
	global.SetMeterProvider(provider)

	b.logger.Debug(Message("prometheus provider started"))

	b.logger.Debug(Message("prometheus metric exporter configured"))

	return provider.Shutdown, nil
}
