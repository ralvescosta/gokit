module github.com/ralvescosta/gokit/logging

go 1.22.2

require (
	github.com/ralvescosta/gokit/configs v1.20.0
	github.com/stretchr/testify v1.9.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ralvescosta/gokit/configs => ../configs
