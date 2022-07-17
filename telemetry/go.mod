module github.com/ralvescosta/gokit/telemetry

go 1.18

require go.opentelemetry.io/otel v1.8.0

require golang.org/x/sys v0.0.0-20220712014510-0a85c31ab51e // indirect

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.7.0
	go.opentelemetry.io/otel/sdk v1.8.0
	go.opentelemetry.io/otel/trace v1.8.0 // indirect
)
