module github.com/ralvescosta/gokit/secrets_manager

go 1.21

require (
	github.com/aws/aws-sdk-go-v2/config v1.25.5
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.23.3
	github.com/ralvescosta/gokit/configs v1.16.0
	github.com/ralvescosta/gokit/logging v1.16.0
	go.uber.org/zap v1.26.0
)

require (
	github.com/aws/aws-sdk-go-v2 v1.23.1 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.16.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.7.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.17.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.20.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.25.4 // indirect
	github.com/aws/smithy-go v1.17.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.1 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ralvescosta/gokit/configs => ../configs

replace github.com/ralvescosta/gokit/logging => ../logging
