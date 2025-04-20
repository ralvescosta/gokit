// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

type MetricsKind string

const (
	OTLP       = MetricsKind("OTLP")
	Prometheus = MetricsKind("Prometheus")
)

type MetricsConfigs struct {
	Enabled bool

	Kind MetricsKind

	OtlpEndpoint string
	OtlpAPIKey   string
}
