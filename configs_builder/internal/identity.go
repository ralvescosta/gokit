package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	"github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadIdentityConfigs() (*configs.IdentityConfigs, error) {
	cfg := &configs.IdentityConfigs{}

	cfg.Domain = os.Getenv(keys.IDENTITY_DOMAIN_ENV_KEY)
	if cfg.Domain == "" {
		return nil, errors.NewErrRequiredConfig(keys.IDENTITY_DOMAIN_ENV_KEY)
	}

	cfg.Audience = os.Getenv(keys.IDENTITY_AUDIENCE_ENV_KEY)
	if cfg.Audience == "" {
		return nil, errors.NewErrRequiredConfig(keys.IDENTITY_AUDIENCE_ENV_KEY)
	}

	cfg.Issuer = os.Getenv(keys.IDENTITY_ISSUER_ENV_KEY)
	if cfg.Issuer == "" {
		return nil, errors.NewErrRequiredConfig(keys.IDENTITY_ISSUER_ENV_KEY)
	}

	cfg.Signature = os.Getenv(keys.IDENTITY_SIGNATURE_ENV_KEY)

	return cfg, nil
}
