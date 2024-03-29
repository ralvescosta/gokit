module github.com/ralvescosta/gokit/sql

go 1.21

require (
	github.com/lib/pq v1.10.9
	github.com/ralvescosta/gokit/configs v1.16.0
	github.com/ralvescosta/gokit/logging v1.16.0
	github.com/stretchr/testify v1.8.4
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.2.3
	go.opentelemetry.io/otel v1.21.0
	go.uber.org/zap v1.26.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.1 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ralvescosta/gokit/configs => ../configs

replace github.com/ralvescosta/gokit/logging => ../logging
