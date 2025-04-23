// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package tracing provides distributed tracing capabilities using OpenTelemetry.
package tracing

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/ralvescosta/gokit/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
)

type (
	// OTLPTracingBuilder extends the TracingBuilder interface with OTLP-specific configuration options
	// such as API key header, connection timeouts, reconnection settings, and compression.
	OTLPTracingBuilder interface {
		TracingBuilder

		// WithAPIKeyHeader adds the API key header from config to the OTLP exporter
		WithAPIKeyHeader() OTLPTracingBuilder

		// WithTimeout sets the timeout duration for OTLP exporter connections
		WithTimeout(t time.Duration) OTLPTracingBuilder

		// WithReconnection sets the reconnection period for the OTLP exporter
		WithReconnection(t time.Duration) OTLPTracingBuilder

		// WithCompression sets the compression algorithm for the OTLP exporter
		WithCompression(c OTLPCompression) OTLPTracingBuilder
	}

	// otlpTracingBuilder is the internal implementation of the OTLPTracingBuilder interface
	otlpTracingBuilder struct {
		tracingBuilder
		reconnectionPeriod time.Duration
		timeout            time.Duration
		compression        OTLPCompression
	}
)

// NewOTLP creates a new instance of OTLPTracingBuilder with default settings from the provided configs
func NewOTLP(cfgs *configs.Configs) OTLPTracingBuilder {
	return &otlpTracingBuilder{
		tracingBuilder: tracingBuilder{
			logger:       cfgs.Logger,
			cfg:          cfgs,
			exporterType: OTLP_TLS_GRPC_EXPORTER,
			headers:      Headers{},
			endpoint:     cfgs.TracingConfigs.OtlpEndpoint,
		},
		reconnectionPeriod: 2 * time.Second,
		timeout:            30 * time.Second,
		compression:        OTLP_GZIP_COMPRESSIONS,
	}
}

// WithAPIKeyHeader adds the API key header from configuration to the OTLP exporter headers
func (b *otlpTracingBuilder) WithAPIKeyHeader() OTLPTracingBuilder {
	b.headers["api-key"] = b.cfg.TracingConfigs.OtlpAPIKey
	return b
}

// AddHeader adds a single header key-value pair to the exporter configuration
func (b *otlpTracingBuilder) AddHeader(key, value string) TracingBuilder {
	b.headers[key] = value
	return b
}

// WithHeaders sets all headers at once for the exporter configuration
func (b *otlpTracingBuilder) WithHeaders(headers Headers) TracingBuilder {
	b.headers = headers
	return b
}

// Type sets the exporter type to be used
func (b *otlpTracingBuilder) Type(t ExporterType) TracingBuilder {
	b.exporterType = t
	return b
}

// Endpoint sets the endpoint URL for the exporter
func (b *otlpTracingBuilder) Endpoint(s string) TracingBuilder {
	b.endpoint = s
	return b
}

// WithTimeout sets the timeout duration for OTLP exporter connections
func (b *otlpTracingBuilder) WithTimeout(t time.Duration) OTLPTracingBuilder {
	b.timeout = t
	return b
}

// WithReconnection sets the reconnection period for the OTLP exporter
func (b *otlpTracingBuilder) WithReconnection(t time.Duration) OTLPTracingBuilder {
	b.reconnectionPeriod = t
	return b
}

// WithCompression sets the compression algorithm for the OTLP exporter
func (b *otlpTracingBuilder) WithCompression(c OTLPCompression) OTLPTracingBuilder {
	b.compression = c
	return b
}

// Build creates and configures the tracing provider based on the builder settings
// Returns a shutdown function to cleanly close the exporter and any error encountered
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

// buildGrpcExporter creates and configures an OTLP gRPC exporter
// Returns a shutdown function to cleanly close the exporter and any error encountered
func (b *otlpTracingBuilder) buildGrpcExporter() (shutdown func(context.Context) error, err error) {
	b.logger.Debug(Message("otlp gRPC trace exporter"))

	var clientOpts = []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(b.endpoint),
		otlptracegrpc.WithReconnectionPeriod(b.reconnectionPeriod),
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
			attribute.String("service.name", b.cfg.AppConfigs.AppName),
			attribute.String("environment", b.cfg.AppConfigs.GoEnv.ToString()),
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
