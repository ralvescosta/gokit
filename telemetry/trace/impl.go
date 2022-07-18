package trace

import (
	"context"
	"errors"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewOTLP(cfg *env.Configs, logger logging.ILogger) TraceBuilder {
	return &traceBuilder{
		logger:             logger,
		cfg:                cfg,
		appName:            cfg.APP_NAME,
		exporterType:       GRPC_EXPORTER,
		endpoint:           cfg.OTLP_ENDPOINT,
		reconnectionPeriod: 2 * time.Second,
		timeout:            30 * time.Second,
		compression:        OTLP_GZIP_COMPRESSIONS,
	}
}

func (b *traceBuilder) WithHeader(headers Headers) TraceBuilder {
	b.headers = headers
	return b
}

func (b *traceBuilder) Type(t OTLPExporterType) TraceBuilder {
	b.exporterType = t
	return b
}

func (b *traceBuilder) Endpoint(s string) TraceBuilder {
	b.endpoint = s
	return b
}

func (b *traceBuilder) WithTimeout(t time.Duration) TraceBuilder {
	b.timeout = t
	return b
}

func (b *traceBuilder) WithReconnection(t time.Duration) TraceBuilder {
	b.reconnectionPeriod = t
	return b
}

func (b *traceBuilder) WithCompression(c OTLPCompression) TraceBuilder {
	b.compression = c
	return b
}

func (b *traceBuilder) Build() (shutdown func(context.Context) error, err error) {
	switch b.exporterType {
	case GRPC_EXPORTER:
		return b.buildTLSGrpcExporter()
	default:
		return nil, errors.New("this pkg support only tls grpc exporter")
	}
}

func (b *traceBuilder) buildTLSGrpcExporter() (shutdown func(context.Context) error, err error) {
	b.logger.Debug("building TLS gRPC Exporter...")

	var clientOpts = []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(b.endpoint),
		otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
		otlptracegrpc.WithReconnectionPeriod(b.reconnectionPeriod),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
		otlptracegrpc.WithTimeout(b.timeout),
		otlptracegrpc.WithHeaders(b.headers),
		otlptracegrpc.WithCompressor(string(b.compression)),
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(clientOpts...),
	)
	if err != nil {
		b.logger.Error("could not create the exporter", logging.ErrorField(err))
		return nil, err
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", b.appName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		b.logger.Error("could not set resources", logging.ErrorField(err))
		return nil, err
	}

	otel.SetTracerProvider(
		sdkTrace.NewTracerProvider(
			sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
			sdkTrace.WithBatcher(exporter),
			sdkTrace.WithResource(resources),
		),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	b.logger.Debug("TLS gRPC Exporter was ...")
	return exporter.Shutdown, nil
}
