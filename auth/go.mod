module github.com/ralvescosta/gokit/auth

go 1.21

require (
	github.com/go-jose/go-jose/v3 v3.0.3
	github.com/ralvescosta/gokit/configs v1.16.0
	go.uber.org/zap v1.26.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.1 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/ralvescosta/gokit/logging v1.16.0
	golang.org/x/crypto v0.19.0 // indirect
)

replace github.com/ralvescosta/gokit/configs => ../configs
