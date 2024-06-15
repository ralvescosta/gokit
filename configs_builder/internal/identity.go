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

	cfg.ClientID = os.Getenv(keys.IdentityClientIDEnvKey)
	if cfg.ClientID == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityClientIDEnvKey)
	}

	cfg.ClientSecret = os.Getenv(keys.IdentityClientSecretEnvKey)
	if cfg.ClientSecret == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityClientSecretEnvKey)
	}

	cfg.GrantType = os.Getenv(keys.IdentityGrantTypeEnvKey)
	if cfg.GrantType == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityGrantTypeEnvKey)
	}

	if milliseconds := os.Getenv(keys.IdentityMillisecondsBetweenJwkEnvKey); milliseconds != "" {
		atoi, err := strconv.Atoi(milliseconds)
		if err != nil {
			return nil, err
		}
		cfg.MillisecondsBetweenJWK = int64(atoi)
	} else {
		// 12 hours
		cfg.MillisecondsBetweenJWK = 43200000
	}

	cfg.Domain = os.Getenv(keys.IdentityDomainEnvKey)
	if cfg.Domain == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityDomainEnvKey)
	}

	cfg.Audience = os.Getenv(keys.IdentityAudienceEnvKey)
	if cfg.Audience == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityAudienceEnvKey)
	}

	cfg.Issuer = os.Getenv(keys.IdentityIssuerEnvKey)
	if cfg.Issuer == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityIssuerEnvKey)
	}

	cfg.Signature = os.Getenv(keys.IdentitySignatureEnvKey)

	return cfg, nil
}
