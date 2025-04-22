# HTTP

## Usage Example

Below is a simple example of how to use the `httpw` package to set up an HTTP server with basic routing and middleware.

```go
package main

import (
    "net/http"
    "github.com/ralvescosta/gokit/httpw/server"
    "github.com/ralvescosta/gokit/logging"
)

func main() {
    // Initialize configurations
    cfgs, err := configsBuilder.
		NewConfigsBuilder().
		RabbitMQ().
		Build()

	if err != nil {
      panic(err)
	}

    logger := cfgs.Logger

    // Create a new HTTP server builder
    serverBuilder := server.NewHTTPServerBuilder(&cfgs)

    // Build the server
    httpServer := serverBuilder.Build()

    // Define a simple handler
    handler := func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    }

    // Register a basic route
    err := httpServer.BasicRoute(http.MethodGet, "/", handler)
    if err != nil {
        logger.Error("Failed to register route", err)
        return
    }

    // Run the server
    if err := httpServer.Run(); err != nil {
        logger.Error("Failed to start server", err)
    }
}
