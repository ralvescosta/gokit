package internal

import (
	"fmt"
	"os"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadHTTPConfigs() (*configs.HTTPConfigs, error) {
	httpConfigs := configs.HTTPConfigs{}

	httpConfigs.Port = os.Getenv(keys.HTTP_PORT_ENV_KEY)
	if httpConfigs.Port == "" {
		return nil, errors.NewErrRequiredConfig(keys.HTTP_PORT_ENV_KEY)
	}

	httpConfigs.Host = os.Getenv(keys.HTTP_HOST_ENV_KEY)
	if httpConfigs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.HTTP_HOST_ENV_KEY)
	}

	httpConfigs.Addr = fmt.Sprintf("%s:%s", httpConfigs.Host, httpConfigs.Port)

	profiling := os.Getenv(keys.HTTP_ENABLE_PROFILING_ENV_KEY)
	if profiling == "true" {
		httpConfigs.EnableProfiling = true
	}

	return &httpConfigs, nil
}
