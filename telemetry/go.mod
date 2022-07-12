module github.com/ralvescosta/toolkit/telemetry

go 1.18

require go.opentelemetry.io/otel v1.7.0

require golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
)
