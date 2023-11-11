package metrics

import (
	"context"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
)

type (
	PrometheusMetricBuilder interface {
		Configs(cfg *configs.Configs) PrometheusMetricBuilder
		Logger(logger logging.Logger) PrometheusMetricBuilder
		HTTPHandler() http.Handler
		Build() (shutdown func(context.Context) error, err error)
	}

	prometheusMetricBuilder struct {
		*basicMetricBuilder
	}
)

func NewPrometheus() PrometheusMetricBuilder {
	return &prometheusMetricBuilder{
		basicMetricBuilder: &basicMetricBuilder{},
	}
}

func (b *prometheusMetricBuilder) Configs(cfg *configs.Configs) PrometheusMetricBuilder {
	b.basicMetricBuilder.cfg = cfg
	b.basicMetricBuilder.appName = cfg.AppConfigs.AppName
	return b
}

func (b *prometheusMetricBuilder) Logger(logger logging.Logger) PrometheusMetricBuilder {
	b.basicMetricBuilder.logger = logger
	return b
}

func (b *prometheusMetricBuilder) HTTPHandler() http.Handler {
	return promhttp.Handler()
}

func (b *prometheusMetricBuilder) Build() (shutdown func(context.Context) error, err error) {
	b.logger.Debug(Message("prometheus metric exporter"))

	b.logger.Debug(Message("creating prometheus resource..."))
	ctx := context.Background()

	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("library.language", "go"),
			attribute.String("service.name", b.appName),
			attribute.String("environment", b.cfg.AppConfigs.GoEnv.ToString()),
			attribute.Int64("ID", int64(os.Getegid())),
		),
	)

	if err != nil {
		b.logger.Error(Message("could not set resources"), zap.Error(err))
		return nil, err
	}
	b.logger.Debug(Message("prometheus resource created"))

	b.logger.Debug(Message("starting prometheus provider..."))

	exporter, err := otelprom.New()
	if err != nil {
		b.logger.Error(Message("error to create prom"), zap.Error(err))
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter), metric.WithResource(resources))

	otel.SetMeterProvider(provider)

	b.logger.Debug(Message("prometheus provider started"))

	b.logger.Debug(Message("prometheus metric exporter configured"))

	return provider.Shutdown, nil
}
