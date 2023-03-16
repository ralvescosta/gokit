module github.com/ralvescosta/gokit/logging

go 1.20

require (
	github.com/ralvescosta/gokit/env v1.0.256
	github.com/stretchr/testify v1.8.2
	go.uber.org/zap v1.24.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ralvescosta/dotenv v1.0.4 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ralvescosta/gokit/env => ../env
