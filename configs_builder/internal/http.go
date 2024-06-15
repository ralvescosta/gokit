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

	httpConfigs.Port = os.Getenv(keys.HTTPPortEnvKey)
	if httpConfigs.Port == "" {
		return nil, errors.NewErrRequiredConfig(keys.HTTPPortEnvKey)
	}

	httpConfigs.Host = os.Getenv(keys.HTTPHostEnvKey)
	if httpConfigs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.HTTPHostEnvKey)
	}

	httpConfigs.Addr = fmt.Sprintf("%s:%s", httpConfigs.Host, httpConfigs.Port)

	profiling := os.Getenv(keys.HTTPEnableProfilingEnvKey)
	if profiling == "true" {
		httpConfigs.EnableProfiling = true
	}

	return &httpConfigs, nil
}
