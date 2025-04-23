// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// MetricsKind represents the type of metrics system to be used in the application.
// It provides a type-safe way to specify different metrics backends.
type MetricsKind string

const (
	// OTLP indicates that metrics should be sent to an OpenTelemetry collector
	OTLP = MetricsKind("OTLP")
	// Prometheus indicates that metrics should be exposed in Prometheus format
	Prometheus = MetricsKind("Prometheus")
)

// MetricsConfigs defines settings for application metrics collection and reporting.
// It provides configuration options for different metrics backends.
type MetricsConfigs struct {
	// Enabled determines whether metrics collection is active
	Enabled bool

	// Kind specifies which metrics system should be used (OTLP or Prometheus)
	Kind MetricsKind

	// OtlpEndpoint specifies the URL of the OpenTelemetry collector endpoint for metrics
	OtlpEndpoint string
	// OtlpAPIKey provides authentication credentials for the OpenTelemetry collector if required
	OtlpAPIKey string
}
