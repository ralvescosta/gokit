package env

import (
	"fmt"
	"os"
)

type (
	HTTPConfigs struct {
		Host string
		Port string
		Addr string
	}
)

func (b *ConfigsBuilderImpl) HTTPServer() ConfigsBuilder {
	b.httpServer = true
	return b
}

func (b *ConfigsBuilderImpl) getHTTPServerConfigs() (*HTTPConfigs, error) {
	if !b.httpServer {
		return nil, nil
	}

	configs := HTTPConfigs{}

	configs.Port = os.Getenv(HTTP_PORT_ENV_KEY)
	if configs.Port == "" {
		return nil, NewErrRequiredConfig(HTTP_PORT_ENV_KEY)
	}

	configs.Host = os.Getenv(HTTP_HOST_ENV_KEY)
	if configs.Host == "" {
		return nil, NewErrRequiredConfig(HTTP_HOST_ENV_KEY)
	}

	configs.Addr = fmt.Sprintf("%s:%s", configs.Host, configs.Port)

	return &configs, nil
}
