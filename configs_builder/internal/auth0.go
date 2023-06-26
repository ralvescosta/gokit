package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadAuth0Configs() (*configs.Auth0Configs, error) {
	cfg := &configs.Auth0Configs{}

	cfg.ClientID = os.Getenv(keys.AUTH0_CLIENT_ID_ENV_KEY)
	if cfg.ClientID == "" {
		return nil, errors.NewErrRequiredConfig(keys.AUTH0_CLIENT_ID_ENV_KEY)
	}

	cfg.ClientSecret = os.Getenv(keys.AUTH0_CLIENT_SECRET_ENV_KEY)
	if cfg.ClientSecret == "" {
		return nil, errors.NewErrRequiredConfig(keys.AUTH0_CLIENT_SECRET_ENV_KEY)
	}

	cfg.GrantType = os.Getenv(keys.AUTH0_GRANT_TYPE_ENV_KEY)
	if cfg.GrantType == "" {
		return nil, errors.NewErrRequiredConfig(keys.AUTH0_GRANT_TYPE_ENV_KEY)
	}

	if milliseconds := os.Getenv(keys.AUTH0_MILLISECONDS_BETWEEN_JWK_ENV_KEY); milliseconds != "" {
		atoi, err := strconv.Atoi(milliseconds)
		if err != nil {
			return nil, err
		}
		cfg.MillisecondsBetweenJWK = int64(atoi)
	} else {
		// 12 hours
		cfg.MillisecondsBetweenJWK = 43200000
	}

	return cfg, nil
}
