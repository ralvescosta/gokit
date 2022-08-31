package metric

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"google.golang.org/grpc"
)

func Otlp() (*otlpmetric.Exporter, error) {
	var clientOpts = []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(""),
		otlpmetricgrpc.WithReconnectionPeriod(time.Minute),
		otlpmetricgrpc.WithDialOption(grpc.WithBlock()),
		otlpmetricgrpc.WithTimeout(time.Minute),
		// otlpmetricgrpc.WithHeaders(b.headers),
		otlpmetricgrpc.WithCompressor("gzip"),
	}

	metricClient := otlpmetricgrpc.NewClient(clientOpts...)

	metric, err := otlpmetric.New(context.Background(), metricClient)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func NewPrometheusExporter() (*prometheus.Exporter, error) {
	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}

	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)

	exporter, err := prometheus.New(config, c)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize prometheus exporter: %w", err)
	}

	global.SetMeterProvider(exporter.MeterProvider())

	return exporter, nil
}
