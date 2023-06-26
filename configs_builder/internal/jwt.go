package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	"github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadJWTConfigs() (*configs.JWTConfigs, error) {
	cfg := &configs.JWTConfigs{}

	cfg.Domain = os.Getenv(keys.JWT_DOMAIN_ENV_KEY)
	if cfg.Domain == "" {
		return nil, errors.NewErrRequiredConfig(keys.JWT_DOMAIN_ENV_KEY)
	}

	cfg.Audience = os.Getenv(keys.JWT_AUDIENCE_ENV_KEY)
	if cfg.Audience == "" {
		return nil, errors.NewErrRequiredConfig(keys.JWT_AUDIENCE_ENV_KEY)
	}

	cfg.Issuer = os.Getenv(keys.JWT_ISSUER_ENV_KEY)
	if cfg.Issuer == "" {
		return nil, errors.NewErrRequiredConfig(keys.JWT_ISSUER_ENV_KEY)
	}

	cfg.Signature = os.Getenv(keys.JWT_SIGNATURE_ENV_KEY)

	return cfg, nil
}
