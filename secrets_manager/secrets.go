// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package secretsmanager provides interfaces and implementations for managing secrets
// across different environments and secret providers.
package secretsmanager

import "context"

type (
	// SecretClient defines the interface for interacting with secret providers.
	// This interface allows for different implementations such as AWS Secrets Manager,
	// Vault, or other secret management solutions.
	SecretClient interface {
		// LoadSecrets loads all secrets from the provider into memory.
		// This method should be called during application initialization.
		LoadSecrets(ctx context.Context) error

		// GetSecret retrieves a specific secret value by its key.
		// Returns the secret value as a string and any error encountered.
		// If the key doesn't exist, an error will be returned.
		GetSecret(ctx context.Context, key string) (string, error)
	}
)
