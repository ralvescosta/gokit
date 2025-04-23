// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	"github.com/ralvescosta/gokit/configs_builder/keys"
)

// ReadIdentityConfigs retrieves authentication and identity-related configuration from environment variables.
// Validates all required identity provider settings and returns an error if any required
// configuration is missing.
func ReadIdentityConfigs() (*configs.IdentityConfigs, error) {
	cfg := &configs.IdentityConfigs{}

	// Get and validate OAuth client ID
	cfg.ClientID = os.Getenv(keys.IdentityClientIDEnvKey)
	if cfg.ClientID == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityClientIDEnvKey)
	}

	// Get and validate OAuth client secret
	cfg.ClientSecret = os.Getenv(keys.IdentityClientSecretEnvKey)
	if cfg.ClientSecret == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityClientSecretEnvKey)
	}

	// Get and validate OAuth grant type
	cfg.GrantType = os.Getenv(keys.IdentityGrantTypeEnvKey)
	if cfg.GrantType == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityGrantTypeEnvKey)
	}

	// Get JWK refresh interval, using default if not specified
	if milliseconds := os.Getenv(keys.IdentityMillisecondsBetweenJwkEnvKey); milliseconds != "" {
		atoi, err := strconv.Atoi(milliseconds)
		if err != nil {
			return nil, err
		}
		cfg.MillisecondsBetweenJWK = int64(atoi)
	} else {
		// Default to 12 hours if not specified
		cfg.MillisecondsBetweenJWK = 43200000
	}

	// Get and validate identity provider domain
	cfg.Domain = os.Getenv(keys.IdentityDomainEnvKey)
	if cfg.Domain == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityDomainEnvKey)
	}

	// Get and validate JWT audience
	cfg.Audience = os.Getenv(keys.IdentityAudienceEnvKey)
	if cfg.Audience == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityAudienceEnvKey)
	}

	// Get and validate JWT issuer
	cfg.Issuer = os.Getenv(keys.IdentityIssuerEnvKey)
	if cfg.Issuer == "" {
		return nil, errors.NewErrRequiredConfig(keys.IdentityIssuerEnvKey)
	}

	// Get JWT signature algorithm (optional)
	cfg.Signature = os.Getenv(keys.IdentitySignatureEnvKey)

	return cfg, nil
}
