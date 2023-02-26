package env

import (
	"fmt"
	"os"
)

const (
	RequiredHTTPServerErrorMessage = "[ConfigBuilder::HTTPServer] %s is required"
)

func (b *ConfigBuilderImpl) HTTPServer() ConfigBuilder {
	b.httpServer = true
	return b
}

func (b *ConfigBuilderImpl) getHTTPServerConfigs() (*HTTPConfigs, error) {
	if !b.httpServer {
		return nil, nil
	}

	configs := HTTPConfigs{}

	configs.Port = os.Getenv(HTTP_PORT_ENV_KEY)
	if configs.Port == "" {
		return nil, fmt.Errorf(RequiredHTTPServerErrorMessage, HTTP_PORT_ENV_KEY)
	}

	configs.Host = os.Getenv(HTTP_HOST_ENV_KEY)
	if configs.Host == "" {
		return nil, fmt.Errorf(RequiredHTTPServerErrorMessage, HTTP_HOST_ENV_KEY)
	}

	configs.Addr = fmt.Sprintf("%s:%s", configs.Host, configs.Port)

	return &configs, nil
}
