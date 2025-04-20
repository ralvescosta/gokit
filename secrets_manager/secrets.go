// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package secretsmanager

import "context"

type (
	SecretClient interface {
		LoadSecrets(ctx context.Context) error
		GetSecret(ctx context.Context, key string) (string, error)
	}
)
