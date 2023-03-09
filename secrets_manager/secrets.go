package secretsmanager

import "context"

type (
	SecretClient interface {
		LoadSecrets(ctx context.Context) error
		GetSecret(ctx context.Context, key string) (string, error)
	}
)
