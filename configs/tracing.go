// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// TracingConfigs defines settings for distributed tracing functionality.
// It allows configuring OpenTelemetry collectors and related tracing options.
type TracingConfigs struct {
	// Enabled determines whether distributed tracing is active for the application
	Enabled bool

	// OtlpEndpoint specifies the URL of the OpenTelemetry collector endpoint
	OtlpEndpoint string
	// OtlpAPIKey provides authentication credentials for the OpenTelemetry collector if required
	OtlpAPIKey string
}
