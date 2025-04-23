# Secrets Manager

The `secretsmanager` package provides a unified interface for managing application secrets across different environments and secret providers.

## Overview

This package offers a clean abstraction for secrets management with the following components:

- `SecretClient` interface - defines the contract for any secrets provider implementation
- AWS Secrets Manager implementation - a ready-to-use client for AWS Secrets Manager

## Usage

### Initializing the AWS Secrets Manager Client

```go
import (
    "github.com/ralvescosta/gokit/configs"
    "github.com/ralvescosta/gokit/secrets_manager"
)

func main() {
    // Initialize configs
    cfgs := configs.New()

    // Create a new AWS Secrets Manager client
    secretClient, err := secretsmanager.NewAwsSecretClient(cfgs)
    if err != nil {
        // Handle error
    }

    // Load secrets during application startup
    err = secretClient.LoadSecrets(context.Background())
    if err != nil {
        // Handle error
    }

    // Use secrets in your application
    apiKey, err := secretClient.GetSecret(context.Background(), "API_KEY")
    if err != nil {
        // Handle error
    }

    // Use apiKey in your application
}
```

## Secret Format

When using the AWS Secrets Manager implementation, secrets should be stored in JSON format with string keys and string values. For example:

```json
{
  "API_KEY": "your-api-key",
  "DATABASE_PASSWORD": "your-db-password",
  "JWT_SECRET": "your-jwt-secret"
}
```

## Secret ID Format

The AWS implementation automatically constructs the secret ID using the following format:

```
{environment}/{application_secret_key}
```

For example, if your application is running in the `development` environment and your application's secret key is `my-app`, the secret ID would be `development/my-app`.

## Benefits

- **Separation of concerns**: Keep your application code separate from your secret management
- **Environment-specific secrets**: Easily manage different secrets for different environments
- **In-memory caching**: Secrets are loaded once and cached in memory for efficient access
- **Extensibility**: Implement the `SecretClient` interface for other secret providers

## Security Considerations

- Secrets are loaded into memory at application startup
- Ensure your application has appropriate IAM permissions to access secrets
- Consider implementing additional encryption for highly sensitive secrets

## Implementing New Secret Providers

To implement a new secret provider, simply implement the `SecretClient` interface:

```go
type MyCustomSecretProvider struct {
    // Your implementation details
}

func (p *MyCustomSecretProvider) LoadSecrets(ctx context.Context) error {
    // Implementation for loading secrets
    return nil
}

func (p *MyCustomSecretProvider) GetSecret(ctx context.Context, key string) (string, error) {
    // Implementation for retrieving a secret
    return "secret-value", nil
}
```
