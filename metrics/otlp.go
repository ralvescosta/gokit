// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

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
	// OTLPMetrics provides a fluent interface to build and configure
	// OpenTelemetry Protocol (OTLP) metrics exporters.
	OTLPMetrics interface {
		// WithAPIKeyHeader adds an API key header to the OTLP exporter.
		WithAPIKeyHeader() OTLPMetrics

		// AddHeader adds a custom header with the specified key and value.
		AddHeader(key, value string) OTLPMetrics

		// WithHeaders sets multiple headers at once for the OTLP exporter.
		WithHeaders(headers Headers) OTLPMetrics

		// Endpoint sets the URL endpoint for the OTLP exporter.
		Endpoint(s string) OTLPMetrics

		// WithTimeout sets the connection timeout duration.
		WithTimeout(t time.Duration) OTLPMetrics

		// WithReconnection sets the reconnection period for the OTLP connection.
		WithReconnection(t time.Duration) OTLPMetrics

		// WithCompression sets the compression type for the OTLP exporter.
		WithCompression(c OTLPCompression) OTLPMetrics

		// Provider creates and configures the OTLP metrics provider.
		// Returns a shutdown function and any error encountered during setup.
		Provider() (shutdown func(context.Context) error, err error)
	}

	// otlpMetrics implements the OTLPMetrics interface.
	otlpMetrics struct {
		*basicMetricsAttr

		// headers contains HTTP headers to be sent with OTLP requests.
		headers Headers

		// endpoint is the URL of the OTLP receiver.
		endpoint string

		// reconnectionPeriod specifies how often to try reconnecting to the endpoint.
		reconnectionPeriod time.Duration

		// timeout specifies the maximum time for OTLP operations.
		timeout time.Duration

		// compression defines the type of compression to use for OTLP data.
		compression OTLPCompression
	}
)

// NewOTLPBuilder creates a new OpenTelemetry metrics builder with default settings.
// It initializes the builder with configuration from the provided configs.
//
// Parameters:
//   - cfgs: Application configuration containing logger and other settings.
//
// Returns:
//   - An OTLPMetrics builder interface for fluent configuration.
func NewOTLPBuilder(cfgs *configs.Configs) OTLPMetrics {
	return &otlpMetrics{
		basicMetricsAttr: &basicMetricsAttr{
			cfg:    cfgs,
			logger: cfgs.Logger,
		},
		reconnectionPeriod: 2 * time.Second,
		timeout:            30 * time.Second,
		compression:        OtlpGzipCompressions,
		headers:            Headers{},
	}
}

// Configs updates the metrics configuration with values from the provided configs.
//
// Parameters:
//   - cfg: The application configuration to use.
//
// Returns:
//   - The OTLPMetrics builder for method chaining.
func (b *otlpMetrics) Configs(cfg *configs.Configs) OTLPMetrics {
	b.basicMetricsAttr.cfg = cfg
	b.basicMetricsAttr.appName = cfg.AppConfigs.AppName
	b.endpoint = cfg.MetricsConfigs.OtlpEndpoint
	return b
}

// Logger sets the logger instance to be used for metrics-related logging.
//
// Parameters:
//   - logger: The logger instance to use.
//
// Returns:
//   - The OTLPMetrics builder for method chaining.
func (b *otlpMetrics) Logger(logger logging.Logger) OTLPMetrics {
	b.basicMetricsAttr.logger = logger
	return b
}

// WithAPIKeyHeader adds an API key header using the key from the application configuration.
//
// Returns:
//   - The OTLPMetrics builder for method chaining.
func (b *otlpMetrics) WithAPIKeyHeader() OTLPMetrics {
	b.headers["api-key"] = b.cfg.MetricsConfigs.OtlpAPIKey
	return b
}

// AddHeader adds a custom header with the provided key and value.
//
// Parameters:
//   - key: The header key.
//   - value: The header value.
//
// Returns:
//   - The OTLPMetrics builder for method chaining.
func (b *otlpMetrics) AddHeader(key, value string) OTLPMetrics {
	b.headers[key] = value
	return b
}

// WithHeaders sets multiple headers at once for the OTLP exporter.
//
// Parameters:
//   - headers: Map of header key-value pairs.
//
// Returns:
//   - The OTLPMetrics builder for method chaining.
func (b *otlpMetrics) WithHeaders(headers Headers) OTLPMetrics {
	b.headers = headers
	return b
}

// Endpoint sets the URL endpoint for the OTLP exporter.
//
// Parameters:
//   - s: The endpoint URL string.
//
// Returns:
//   - The OTLPMetrics builder for method chaining.
func (b *otlpMetrics) Endpoint(s string) OTLPMetrics {
	b.endpoint = s
	return b
}

// WithTimeout sets the connection timeout duration.
//
// Parameters:
//   - t: The timeout duration.
//
// Returns:
//   - The OTLPMetrics builder for method chaining.
func (b *otlpMetrics) WithTimeout(t time.Duration) OTLPMetrics {
	b.timeout = t
	return b
}

// WithReconnection sets the reconnection period for the OTLP connection.
//
// Parameters:
//   - t: The reconnection period duration.
//
// Returns:
//   - The OTLPMetrics builder for method chaining.
func (b *otlpMetrics) WithReconnection(t time.Duration) OTLPMetrics {
	b.reconnectionPeriod = t
	return b
}

// WithCompression sets the compression type for the OTLP exporter.
//
// Parameters:
//   - c: The compression type to use.
//
// Returns:
//   - The OTLPMetrics builder for method chaining.
func (b *otlpMetrics) WithCompression(c OTLPCompression) OTLPMetrics {
	b.compression = c
	return b
}

// Provider creates and configures the OTLP metrics provider.
// It sets up the exporter, resource attributes, meter provider, and registers
// everything with the OpenTelemetry global context.
//
// Returns:
//   - shutdown: A function to properly shut down the metrics exporter.
//   - err: Any error encountered during setup.
func (b *otlpMetrics) Provider() (shutdown func(context.Context) error, err error) {
	b.logger.Debug(Message("otlp gRPC metric exporter"))

	var clientOpts = []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(b.endpoint),
		otlpmetricgrpc.WithReconnectionPeriod(b.reconnectionPeriod),
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
