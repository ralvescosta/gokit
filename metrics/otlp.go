package metrics

import (
	"context"
	"os"
	"time"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
)

type (
	OTLPMetrics interface {
		WithAPIKeyHeader() OTLPMetrics
		AddHeader(key, value string) OTLPMetrics
		WithHeaders(headers Headers) OTLPMetrics
		Endpoint(s string) OTLPMetrics
		WithTimeout(t time.Duration) OTLPMetrics
		WithReconnection(t time.Duration) OTLPMetrics
		WithCompression(c OTLPCompression) OTLPMetrics
		Provider() (shutdown func(context.Context) error, err error)
	}

	otlpMetrics struct {
		*basicMetricsAttr

		headers            Headers
		endpoint           string
		reconnectionPeriod time.Duration
		timeout            time.Duration
		compression        OTLPCompression
	}
)

func NewOTLPBuilder(cfg *configs.Configs, logger logging.Logger) OTLPMetrics {
	return &otlpMetrics{
		basicMetricsAttr: &basicMetricsAttr{
			cfg:    cfg,
			logger: logger,
		},
		reconnectionPeriod: 2 * time.Second,
		timeout:            30 * time.Second,
		compression:        OTLP_GZIP_COMPRESSIONS,
		headers:            Headers{},
	}
}

func (b *otlpMetrics) Configs(cfg *configs.Configs) OTLPMetrics {
	b.basicMetricsAttr.cfg = cfg
	b.basicMetricsAttr.appName = cfg.AppConfigs.AppName
	b.endpoint = cfg.MetricsConfigs.OtlpEndpoint
	return b
}

func (b *otlpMetrics) Logger(logger logging.Logger) OTLPMetrics {
	b.basicMetricsAttr.logger = logger
	return b
}

func (b *otlpMetrics) WithAPIKeyHeader() OTLPMetrics {
	b.headers["api-key"] = b.cfg.MetricsConfigs.OtlpAPIKey
	return b
}

func (b *otlpMetrics) AddHeader(key, value string) OTLPMetrics {
	b.headers[key] = value
	return b
}

func (b *otlpMetrics) WithHeaders(headers Headers) OTLPMetrics {
	b.headers = headers
	return b
}

func (b *otlpMetrics) Endpoint(s string) OTLPMetrics {
	b.endpoint = s
	return b
}

func (b *otlpMetrics) WithTimeout(t time.Duration) OTLPMetrics {
	b.timeout = t
	return b
}

func (b *otlpMetrics) WithReconnection(t time.Duration) OTLPMetrics {
	b.reconnectionPeriod = t
	return b
}

func (b *otlpMetrics) WithCompression(c OTLPCompression) OTLPMetrics {
	b.compression = c
	return b
}

func (b *otlpMetrics) Provider() (shutdown func(context.Context) error, err error) {
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
			grpc.WithUserAgent("OTel OTLP Exporter Go/"+otel.Version()),
		),
		otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
	}

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
			attribute.String("environment", b.cfg.AppConfigs.GoEnv.ToString()),
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
			metric.NewPeriodicReader(exporter, metric.WithInterval(5*time.Second), metric.WithTimeout(10*time.Second)),
		),
		metric.WithResource(resources),
	)
	b.logger.Debug(Message("otlp provider was configured"))

	otel.SetMeterProvider(provider)

	b.logger.Debug(Message("otlp gRPC metric exporter configured"))
	return exporter.Shutdown, nil
}
