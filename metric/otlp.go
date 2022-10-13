package metric

import (
	"context"
	"os"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
)

func NewOTLP(cfg *env.Config, logger logging.Logger) OTLPMetricBuilder {
	return &otlpMetricBuilder{
		basicMetricBuilder: basicMetricBuilder{
			logger:  logger,
			cfg:     cfg,
			appName: cfg.APP_NAME,
		},
		endpoint:           cfg.OTLP_ENDPOINT,
		reconnectionPeriod: 2 * time.Second,
		timeout:            30 * time.Second,
		compression:        OTLP_GZIP_COMPRESSIONS,
		headers:            Headers{},
	}
}

func (b *otlpMetricBuilder) WithApiKeyHeader() OTLPMetricBuilder {
	b.headers["api-key"] = b.cfg.OTLP_API_KEY
	return b
}

func (b *otlpMetricBuilder) AddHeader(key, value string) OTLPMetricBuilder {
	b.headers[key] = value
	return b
}

func (b *otlpMetricBuilder) WithHeaders(headers Headers) OTLPMetricBuilder {
	b.headers = headers
	return b
}

func (b *otlpMetricBuilder) Endpoint(s string) OTLPMetricBuilder {
	b.endpoint = s
	return b
}

func (b *otlpMetricBuilder) WithTimeout(t time.Duration) OTLPMetricBuilder {
	b.timeout = t
	return b
}

func (b *otlpMetricBuilder) WithReconnection(t time.Duration) OTLPMetricBuilder {
	b.reconnectionPeriod = t
	return b
}

func (b *otlpMetricBuilder) WithCompression(c OTLPCompression) OTLPMetricBuilder {
	b.compression = c
	return b
}

func (b *otlpMetricBuilder) Build() (shutdown func(context.Context) error, err error) {
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
	ctx := context.Background()

	b.logger.Debug(Message("connecting to otlp exporter..."))
	exporter, err := otlpmetricgrpc.New(ctx, clientOpts...)
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
	provider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(exporter, metric.WithInterval(2*time.Second), metric.WithTimeout(10*time.Second)),
		),
		metric.WithResource(resources),
	)
	b.logger.Debug(Message("otlp provider was configured"))

	global.SetMeterProvider(provider)

	b.logger.Debug(Message("otlp gRPC metric exporter configured"))
	return exporter.Shutdown, nil
}
