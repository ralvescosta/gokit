package metric

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
)

func NewOTLP(cfg *env.Config, logger logging.ILogger) MetricBuilder {
	return &metricBuilder{
		logger:             logger,
		cfg:                cfg,
		appName:            cfg.APP_NAME,
		exporterType:       OTLP_GRPC_TLS_EXPORTER,
		endpoint:           cfg.OTLP_ENDPOINT,
		reconnectionPeriod: 2 * time.Second,
		timeout:            30 * time.Second,
		compression:        OTLP_GZIP_COMPRESSIONS,
		headers:            Headers{},
	}
}

func (b *metricBuilder) WithApiKeyHeader() MetricBuilder {
	b.headers["api-key"] = b.cfg.OTLP_API_KEY
	return b
}

func (b *metricBuilder) AddHeader(key, value string) MetricBuilder {
	b.headers[key] = value
	return b
}

func (b *metricBuilder) WithHeaders(headers Headers) MetricBuilder {
	b.headers = headers
	return b
}

func (b *metricBuilder) Type(t MetricExporterType) MetricBuilder {
	b.exporterType = t
	return b
}

func (b *metricBuilder) Endpoint(s string) MetricBuilder {
	b.endpoint = s
	return b
}

func (b *metricBuilder) WithTimeout(t time.Duration) MetricBuilder {
	b.timeout = t
	return b
}

func (b *metricBuilder) WithReconnection(t time.Duration) MetricBuilder {
	b.reconnectionPeriod = t
	return b
}

func (b *metricBuilder) WithCompression(c OTLPCompression) MetricBuilder {
	b.compression = c
	return b
}

func (b *metricBuilder) Build(ctx context.Context) (shutdown func(context.Context) error, err error) {
	switch b.exporterType {
	case OTLP_GRPC_EXPORTER:
		fallthrough
	case OTLP_GRPC_TLS_EXPORTER:
		return b.otlpGrpcExporter(ctx)
	case PROMETHEUS_EXPORTER:
		return b.prometheusExporter(ctx)
	default:
		return nil, errors.New("")
	}
}

func (b *metricBuilder) otlpGrpcExporter(ctx context.Context) (shutdown func(context.Context) error, err error) {
	b.logger.Debug(LogMessage("otlp gRPC metric exporter"))

	var clientOpts = []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(b.endpoint),
		otlpmetricgrpc.WithReconnectionPeriod(b.reconnectionPeriod),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
		otlpmetricgrpc.WithTimeout(b.timeout),
		otlpmetricgrpc.WithHeaders(b.headers),
		otlpmetricgrpc.WithCompressor(string(b.compression)),
		otlptracegrpc.WithDialOption(
			grpc.WithConnectParams(grpc.ConnectParams{
				Backoff: backoff.Config{
					BaseDelay:  1 * time.Second,
					Multiplier: 1.6,
					MaxDelay:   15 * time.Second,
				},
				MinConnectTimeout: 0,
			}),
		),
	}

	if b.exporterType == OTLP_GRPC_TLS_EXPORTER {
		clientOpts = append(clientOpts, otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	} else {
		clientOpts = append(clientOpts, otlpmetricgrpc.WithInsecure())
	}

	b.logger.Debug(LogMessage("connecting to otlp exporter..."))
	exporter, err := otlpmetric.New(ctx, otlpmetricgrpc.NewClient(clientOpts...))
	if err != nil {
		b.logger.Error(LogMessage("could not create the exporter"), logging.ErrorField(err))
		return nil, err
	}
	b.logger.Debug(LogMessage("otlp exporter connected"))

	b.logger.Debug(LogMessage("creating otlp resource..."))
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
	b.logger.Debug(LogMessage("otlp resource created"))

	b.logger.Debug(LogMessage("configure otlp provider..."))
	metricProvider := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(2*time.Second),
		controller.WithResource(resources),
	)
	b.logger.Debug(LogMessage("otlp provider was configured"))

	global.SetMeterProvider(metricProvider)

	b.logger.Debug(LogMessage("starting otlp provider..."))
	if err := metricProvider.Start(ctx); err != nil {
		b.logger.Error(LogMessage("could not started the provider"), logging.ErrorField(err))
		return nil, err
	}
	b.logger.Debug(LogMessage("otlp provider started"))

	b.logger.Debug(LogMessage("otlp gRPC metric exporter configured"))
	return exporter.Shutdown, nil
}

func (b *metricBuilder) prometheusExporter(ctx context.Context) (shutdown func(context.Context) error, err error) {
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
