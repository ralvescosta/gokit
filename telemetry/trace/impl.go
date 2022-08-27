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
		exporterType:       TLS_GRPC_EXPORTER,
		endpoint:           cfg.OTLP_ENDPOINT,
		reconnectionPeriod: 2 * time.Second,
		timeout:            30 * time.Second,
		compression:        OTLP_GZIP_COMPRESSIONS,
		headers:            Headers{},
	}
}

func (b *traceBuilder) WithApiKeyHeader() TraceBuilder {
	b.headers["api-key"] = b.cfg.OTLP_API_KEY
	return b
}

func (b *traceBuilder) AddHeader(key, value string) TraceBuilder {
	b.headers[key] = value
	return b
}

func (b *traceBuilder) WithHeaders(headers Headers) TraceBuilder {
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

func (b *traceBuilder) Build(ctx context.Context) (shutdown func(context.Context) error, err error) {
	switch b.exporterType {
	case GRPC_EXPORTER:
		fallthrough
	case TLS_GRPC_EXPORTER:
		return b.buildGrpcExporter(ctx)
	default:
		return nil, errors.New("this pkg support only grpc exporter")
	}
}

func (b *traceBuilder) buildGrpcExporter(ctx context.Context) (shutdown func(context.Context) error, err error) {
	b.logger.Debug(LogMessage("otlp gRPC trace exporter"))

	var clientOpts = []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(b.endpoint),
		otlptracegrpc.WithReconnectionPeriod(b.reconnectionPeriod),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
		otlptracegrpc.WithTimeout(b.timeout),
		otlptracegrpc.WithHeaders(b.headers),
		otlptracegrpc.WithCompressor(string(b.compression)),
	}

	if b.exporterType == TLS_GRPC_EXPORTER {
		clientOpts = append(clientOpts, otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	} else {
		clientOpts = append(clientOpts, otlptracegrpc.WithInsecure())
	}

	b.logger.Debug(LogMessage("connecting to otlp exporter..."))
	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(clientOpts...),
	)
	if err != nil {
		b.logger.Error(LogMessage("could not create the exporter"), logging.ErrorField(err))
		return nil, err
	}

	b.logger.Debug(LogMessage("creating otlp resource..."))
	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", b.appName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		b.logger.Error(LogMessage("could not set resources"), logging.ErrorField(err))
		return nil, err
	}

	b.logger.Debug(LogMessage("setting otlp provider..."))
	otel.SetTracerProvider(
		sdkTrace.NewTracerProvider(
			sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
			sdkTrace.WithBatcher(exporter),
			sdkTrace.WithResource(resources),
		),
	)

	b.logger.Debug(LogMessage("setting otlp propagator..."))
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	b.logger.Debug(LogMessage("tls grpc exporter was configured"))
	return exporter.Shutdown, nil
}
