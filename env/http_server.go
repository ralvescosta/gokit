package env

import (
	"fmt"
	"os"
)

const (
	RequiredHTTPServerErrorMessage = "[ConfigBuilder::HTTPServer] %s is required"
)

func (c *Configs) HTTPServer() IConfigs {
	if c.Err != nil {
		return c
	}

	c.HTTP_PORT = os.Getenv(HTTP_PORT_ENV_KEY)
	if c.HTTP_PORT == "" {
		c.Err = fmt.Errorf(RequiredHTTPServerErrorMessage, HTTP_PORT_ENV_KEY)
		return c
	}

	c.HTTP_HOST = os.Getenv(HTTP_HOST_ENV_KEY)
	if c.HTTP_HOST == "" {
		c.Err = fmt.Errorf(RequiredHTTPServerErrorMessage, HTTP_HOST_ENV_KEY)
		return c
	}

	c.HTTP_ADDR = fmt.Sprintf("%s:%s", c.HTTP_HOST, c.HTTP_PORT)

	return c
}
