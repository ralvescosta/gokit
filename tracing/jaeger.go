package tracing

import (
	"context"
	"errors"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func NewJaeger(cfg *env.Config, logger logging.ILogger) JaegerTracingBuilder {
	return &jaegerTracingBuilder{
		logger:       logger,
		cfg:          cfg,
		appName:      cfg.APP_NAME,
		exporterType: JAEGER_EXPORTER,
		endpoint:     cfg.JAEGER_AGENT_HOST,
		headers:      Headers{},
	}
}

func (b *jaegerTracingBuilder) AddHeader(key, value string) TracingBuilder {
	b.headers[key] = value
	return b
}

func (b *jaegerTracingBuilder) WithHeaders(headers Headers) TracingBuilder {
	b.headers = headers
	return b
}

func (b *jaegerTracingBuilder) Type(t ExporterType) TracingBuilder {
	b.exporterType = t
	return b
}

func (b *jaegerTracingBuilder) Endpoint(s string) TracingBuilder {
	b.endpoint = s
	return b
}

func (b *jaegerTracingBuilder) Build(ctx context.Context) (shutdown func(context.Context) error, err error) {
	switch b.exporterType {
	case JAEGER_EXPORTER:
		return b.buildJaegerExporter(ctx)
	default:
		return nil, errors.New("this pkg support only grpc exporter")
	}
}

func (b *jaegerTracingBuilder) buildJaegerExporter(ctx context.Context) (shutdown func(context.Context) error, err error) {
	b.logger.Debug(LogMessage("jaeger trace exporter"))
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(b.endpoint)))
	if err != nil {
		return nil, err
	}

	b.logger.Debug(LogMessage("configuring jaeger provider..."))
	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exp),
		sdkTrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(b.appName),
			attribute.String("environment", b.cfg.GO_ENV.ToString()),
			attribute.Int64("ID", 1999),
		)),
	)
	b.logger.Debug(LogMessage("jaeger provider configured"))

	b.logger.Debug(LogMessage("configuring jaeger propagator..."))
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	b.logger.Debug(LogMessage("propagator configured"))

	b.logger.Debug(LogMessage("configure jaeger as a default exporter"))
	otel.SetTracerProvider(tp)
	b.logger.Debug(LogMessage("default exporter configured"))

	b.logger.Debug(LogMessage("jaeger trace exporter configured"))
	return tp.Shutdown, nil
}