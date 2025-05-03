module github.com/ralvescosta/gokit/configs_builder

go 1.24.0

require github.com/ralvescosta/gokit/configs v1.21.0

require (
	github.com/joho/godotenv v1.5.1
	github.com/ralvescosta/gokit/logging v1.32.0
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

replace github.com/ralvescosta/gokit/configs => ../configs
