package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	"github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadIdentityConfigs() (*configs.IdentityConfigs, error) {
	cfg := &configs.IdentityConfigs{}

	cfg.ClientID = os.Getenv(keys.IDENTITY_CLIENT_ID_ENV_KEY)
	if cfg.ClientID == "" {
		return nil, errors.NewErrRequiredConfig(keys.IDENTITY_CLIENT_ID_ENV_KEY)
	}

	cfg.ClientSecret = os.Getenv(keys.IDENTITY_CLIENT_SECRET_ENV_KEY)
	if cfg.ClientSecret == "" {
		return nil, errors.NewErrRequiredConfig(keys.IDENTITY_CLIENT_SECRET_ENV_KEY)
	}

	cfg.GrantType = os.Getenv(keys.IDENTITY_GRANT_TYPE_ENV_KEY)
	if cfg.GrantType == "" {
		return nil, errors.NewErrRequiredConfig(keys.IDENTITY_GRANT_TYPE_ENV_KEY)
	}

	if milliseconds := os.Getenv(keys.IDENTITY_MILLISECONDS_BETWEEN_JWK_ENV_KEY); milliseconds != "" {
		atoi, err := strconv.Atoi(milliseconds)
		if err != nil {
			return nil, err
		}
		cfg.MillisecondsBetweenJWK = int64(atoi)
	} else {
		// 12 hours
		cfg.MillisecondsBetweenJWK = 43200000
	}

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
