package tracing

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
)

func NewOTLP(cfg *env.Config, logger logging.Logger) OTLPTracingBuilder {
	return &otlpTracingBuilder{
		tracingBuilder: tracingBuilder{
			logger:       logger,
			cfg:          cfg,
			appName:      cfg.APP_NAME,
			exporterType: OTLP_TLS_GRPC_EXPORTER,
			endpoint:     cfg.OTLP_ENDPOINT,
			headers:      Headers{},
		},
		reconnectionPeriod: 2 * time.Second,
		timeout:            30 * time.Second,
		compression:        OTLP_GZIP_COMPRESSIONS,
	}
}

func (b *otlpTracingBuilder) WithApiKeyHeader() OTLPTracingBuilder {
	b.headers["api-key"] = b.cfg.OTLP_API_KEY
	return b
}

func (b *otlpTracingBuilder) AddHeader(key, value string) TracingBuilder {
	b.headers[key] = value
	return b
}

func (b *otlpTracingBuilder) WithHeaders(headers Headers) TracingBuilder {
	b.headers = headers
	return b
}

func (b *otlpTracingBuilder) Type(t ExporterType) TracingBuilder {
	b.exporterType = t
	return b
}

func (b *otlpTracingBuilder) Endpoint(s string) TracingBuilder {
	b.endpoint = s
	return b
}

func (b *otlpTracingBuilder) WithTimeout(t time.Duration) OTLPTracingBuilder {
	b.timeout = t
	return b
}

func (b *otlpTracingBuilder) WithReconnection(t time.Duration) OTLPTracingBuilder {
	b.reconnectionPeriod = t
	return b
}

func (b *otlpTracingBuilder) WithCompression(c OTLPCompression) OTLPTracingBuilder {
	b.compression = c
	return b
}

func (b *otlpTracingBuilder) Build() (shutdown func(context.Context) error, err error) {
	switch b.exporterType {
	case OTLP_GRPC_EXPORTER:
		fallthrough
	case OTLP_TLS_GRPC_EXPORTER:
		return b.buildGrpcExporter()
	default:
		return nil, errors.New("this pkg support only grpc exporter")
	}
}

func (b *otlpTracingBuilder) buildGrpcExporter() (shutdown func(context.Context) error, err error) {
	b.logger.Debug(Message("otlp gRPC trace exporter"))

	var clientOpts = []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(b.endpoint),
		otlptracegrpc.WithReconnectionPeriod(b.reconnectionPeriod),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
		otlptracegrpc.WithTimeout(b.timeout),
		otlptracegrpc.WithHeaders(b.headers),
		otlptracegrpc.WithCompressor(string(b.compression)),
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

	if b.exporterType == OTLP_TLS_GRPC_EXPORTER {
		clientOpts = append(clientOpts, otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	} else {
		clientOpts = append(clientOpts, otlptracegrpc.WithInsecure())
	}

	b.logger.Debug(Message("connecting to otlp exporter..."))
	ctx := context.Background()
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(clientOpts...),
	)
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

	b.logger.Debug(Message("configuring otlp provider..."))
	otel.SetTracerProvider(
		sdkTrace.NewTracerProvider(
			sdkTrace.WithSampler(sdkTrace.ParentBased(sdkTrace.TraceIDRatioBased(0.85))),
			sdkTrace.WithBatcher(exporter),
			sdkTrace.WithResource(resources),
		),
	)
	b.logger.Debug(Message("otlp provider configured"))

	b.logger.Debug(Message("configuring otlp propagator..."))
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	b.logger.Debug(Message("tls grpc exporter was configured"))

	b.logger.Debug(Message("otlp gRPC trace exporter configured"))
	return exporter.Shutdown, nil
}
