module github.com/ralvescosta/gokit/auth

go 1.24.0

require (
	github.com/go-jose/go-jose/v3 v3.0.4
	github.com/ralvescosta/gokit/configs v1.21.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/ralvescosta/gokit/logging v1.21.0
	golang.org/x/crypto v0.37.0 // indirect
)

replace github.com/ralvescosta/gokit/configs => ../configs

replace github.com/ralvescosta/gokit/logging => ../logging
