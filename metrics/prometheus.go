// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package metrics

import (
	"context"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ralvescosta/gokit/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/zap"
)

type (
	// PrometheusMetrics defines the interface for Prometheus metrics exporters.
	PrometheusMetrics interface {
		// HTTPHandler returns an HTTP handler that exposes Prometheus metrics.
		HTTPHandler() http.Handler

		// Provider creates and configures the Prometheus metrics provider.
		// Returns a shutdown function and any error encountered during setup.
		Provider() (shutdown func(context.Context) error, err error)
	}

	// prometheusMetrics implements the PrometheusMetrics interface.
	prometheusMetrics struct {
		*basicMetricsAttr
	}
)

// NewPrometheus creates a new Prometheus metrics exporter with the given configuration.
//
// Parameters:
//   - cfgs: Application configuration containing logger and other settings.
//
// Returns:
//   - A PrometheusMetrics interface for configuring and using Prometheus metrics.
func NewPrometheus(cfgs *configs.Configs) PrometheusMetrics {
	return &prometheusMetrics{
		basicMetricsAttr: &basicMetricsAttr{cfg: cfgs, logger: cfgs.Logger},
	}
}

// HTTPHandler returns an HTTP handler that can be mounted to expose Prometheus
// metrics for scraping.
//
// Returns:
//   - An http.Handler that serves Prometheus metrics.
func (b *prometheusMetrics) HTTPHandler() http.Handler {
	return promhttp.Handler()
}

// Provider creates and configures the Prometheus metrics provider.
// It sets up the exporter, resource attributes, meter provider, and registers
// everything with the OpenTelemetry global context.
//
// Returns:
//   - shutdown: A function to properly shut down the metrics provider.
//   - err: Any error encountered during setup.
func (b *prometheusMetrics) Provider() (shutdown func(context.Context) error, err error) {
	b.logger.Debug(Message("prometheus metric exporter"))

	b.logger.Debug(Message("creating prometheus resource..."))
	ctx := context.Background()

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
	b.logger.Debug(Message("prometheus resource created"))

	b.logger.Debug(Message("starting prometheus provider..."))

	exporter, err := otelprom.New()
	if err != nil {
		b.logger.Error(Message("error to create prom"), zap.Error(err))
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter), metric.WithResource(resources))

	otel.SetMeterProvider(provider)

	b.logger.Debug(Message("prometheus provider started"))

	b.logger.Debug(Message("prometheus metric exporter configured"))

	return provider.Shutdown, nil
}
