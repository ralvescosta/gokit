package tracing

import (
	"context"
	"errors"
	"os"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type (
	JaegerTracingBuilder interface {
		TracingBuilder
	}

	jaegerTracingBuilder struct {
		tracingBuilder
	}
)

func NewJaeger(cfg *configs.Configs, logger logging.Logger) JaegerTracingBuilder {

	return &jaegerTracingBuilder{
		tracingBuilder: tracingBuilder{
			logger:       logger,
			cfg:          cfg,
			exporterType: JAEGER_EXPORTER,
			headers:      Headers{},
		},
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

func (b *jaegerTracingBuilder) Build() (shutdown func(context.Context) error, err error) {
	switch b.exporterType {
	case JAEGER_EXPORTER:
		return b.buildJaegerExporter()
	default:
		return nil, errors.New("this pkg support only grpc exporter")
	}
}

func (b *jaegerTracingBuilder) buildJaegerExporter() (shutdown func(context.Context) error, err error) {
	b.logger.Debug(Message("jaeger trace exporter"))
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(b.endpoint)))
	if err != nil {
		return nil, err
	}

	b.logger.Debug(Message("configuring jaeger provider..."))
	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exp),
		sdkTrace.WithSampler(
			sdkTrace.ParentBased(
				sdkTrace.TraceIDRatioBased(0.01),
			),
		),
		sdkTrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			attribute.String("library.language", "go"),
			semconv.ServiceNameKey.String(b.cfg.AppConfigs.AppName),
			attribute.String("environment", b.cfg.AppConfigs.GoEnv.ToString()),
			attribute.Int64("ID", int64(os.Getegid())),
		)),
	)
	b.logger.Debug(Message("jaeger provider configured"))

	b.logger.Debug(Message("configuring jaeger propagator..."))
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	b.logger.Debug(Message("propagator configured"))

	b.logger.Debug(Message("configure jaeger as a default exporter"))
	otel.SetTracerProvider(tp)
	b.logger.Debug(Message("default exporter configured"))

	b.logger.Debug(Message("jaeger trace exporter configured"))
	return tp.Shutdown, nil
}
