package metric

import (
	"context"
	"os"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
)

func NewOTLP(cfg *env.Config, logger logging.Logger) MetricBuilder {
	return &otlpMetricBuilder{
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

func (b *otlpMetricBuilder) WithApiKeyHeader() MetricBuilder {
	b.headers["api-key"] = b.cfg.OTLP_API_KEY
	return b
}

func (b *otlpMetricBuilder) AddHeader(key, value string) MetricBuilder {
	b.headers[key] = value
	return b
}

func (b *otlpMetricBuilder) WithHeaders(headers Headers) MetricBuilder {
	b.headers = headers
	return b
}

func (b *otlpMetricBuilder) Endpoint(s string) MetricBuilder {
	b.endpoint = s
	return b
}

func (b *otlpMetricBuilder) WithTimeout(t time.Duration) MetricBuilder {
	b.timeout = t
	return b
}

func (b *otlpMetricBuilder) WithReconnection(t time.Duration) MetricBuilder {
	b.reconnectionPeriod = t
	return b
}

func (b *otlpMetricBuilder) WithCompression(c OTLPCompression) MetricBuilder {
	b.compression = c
	return b
}

func (b *otlpMetricBuilder) Build(ctx context.Context) (shutdown func(context.Context) error, err error) {
	return b.otlpGrpcExporter(ctx)
}

func (b *otlpMetricBuilder) otlpGrpcExporter(ctx context.Context) (shutdown func(context.Context) error, err error) {
	b.logger.Debug(Message("otlp gRPC metric exporter"))

	var clientOpts = []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(b.endpoint),
		otlpmetricgrpc.WithReconnectionPeriod(b.reconnectionPeriod),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
		otlpmetricgrpc.WithTimeout(b.timeout),
		otlpmetricgrpc.WithHeaders(b.headers),
		otlpmetricgrpc.WithCompressor(string(b.compression)),
		otlpmetricgrpc.WithDialOption(
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

	clientOpts = append(clientOpts, otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))

	b.logger.Debug(Message("connecting to otlp exporter..."))
	exporter, err := otlpmetric.New(ctx, otlpmetricgrpc.NewClient(clientOpts...))
	if err != nil {
		b.logger.Error(Message("could not create the exporter"), zap.Error(err))
		return nil, err
	}
	b.logger.Debug(Message("otlp exporter connected"))

	b.logger.Debug(Message("creating otlp resource..."))
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
	b.logger.Debug(Message("otlp resource created"))

	b.logger.Debug(Message("configure otlp provider..."))
	metricProvider := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(2*time.Second),
		controller.WithResource(resources),
	)
	b.logger.Debug(Message("otlp provider was configured"))

	global.SetMeterProvider(metricProvider)

	b.logger.Debug(Message("starting otlp provider..."))
	if err := metricProvider.Start(ctx); err != nil {
		b.logger.Error(Message("could not started the provider"), zap.Error(err))
		return nil, err
	}
	b.logger.Debug(Message("otlp provider started"))

	b.logger.Debug(Message("otlp gRPC metric exporter configured"))
	return exporter.Shutdown, nil
}
