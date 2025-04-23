# HTTP Wrapper (httpw) Package

[![GoDoc](https://godoc.org/github.com/ralvescosta/gokit/httpw?status.svg)](https://godoc.org/github.com/ralvescosta/gokit/httpw)
[![Go Report Card](https://goreportcard.com/badge/github.com/ralvescosta/gokit/httpw)](https://goreportcard.com/report/github.com/ralvescosta/gokit/httpw)

The `httpw` package provides a comprehensive set of utilities for building robust HTTP services in Go. It offers a flexible API for creating and managing HTTP servers, defining routes, registering middleware, handling request validation, and formatting responses.

## Features

- **HTTP Server Management**: Create and configure HTTP servers with support for graceful shutdown
- **Routing**: Define routes with fluent builder patterns
- **Middleware**: Apply cross-cutting concerns like authentication
- **Request Validation**: Validate request bodies against structural rules
- **Response Formatting**: Build standardized HTTP responses
- **Error Handling**: Create consistent error responses

## Installation

```bash
go get github.com/ralvescosta/gokit/httpw
```

## Quick Start

Here's a simple example of how to create an HTTP server using the `httpw` package:

```go
package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder"
	"github.com/ralvescosta/gokit/httpw/server"
	"github.com/ralvescosta/gokit/httpw/viewmodels"
)

func main() {
	// Create configuration
	cfg := configs_builder.NewConfigsBuilder().Build()
	
	// Create a signal channel for graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Create an HTTP server
	httpServer := server.NewHTTPServerBuilder(cfg).
		WithTracing().
		WithMetrics().
		Signal(sig).
		Build()

	// Define a route handler
	helloHandler := func(w http.ResponseWriter, r *http.Request) {
		viewmodels.NewResponseBuilder().
			Writer(w).
			Ok().
			JSON(map[string]string{"message": "Hello, World!"}).
			Build()
	}

	// Register a route
	httpServer.BasicRoute(http.MethodGet, "/hello", helloHandler)

	// Start the server
	if err := httpServer.Run(); err != nil {
		panic(err)
	}
}
```

## Package Structure

- **httpw**: Root package with utility functions and error definitions
- **server**: Components for building and managing HTTP servers
- **middlewares**: HTTP middleware components (e.g., authentication)
- **validator**: Request validation utilities
- **viewmodels**: Standardized response structures and builders

## Server Configuration

The HTTP server can be configured with various options using the builder pattern:

```go
httpServer := server.NewHTTPServerBuilder(cfg).
    WithTLS().                          // Enable TLS
    WithTracing().                      // Enable OpenTelemetry tracing
    WithMetrics().                      // Enable metrics collection
    WithOpenAPI().                      // Enable OpenAPI documentation
    ExportPrometheusScraping().         // Enable Prometheus metrics endpoint
    Timeouts(5*time.Second,             // Configure custom timeouts
             10*time.Second,
             30*time.Second).
    Signal(sig).                        // Set signal channel for graceful shutdown
    Build()
```

## Route Registration

Routes can be registered using different approaches:

### Basic Route

```go
httpServer.BasicRoute(http.MethodGet, "/hello", helloHandler)
```

### Route Builder

```go
route := server.NewRouteBuilder().
    GET("/users").
    Handler(getUsersHandler).
    Middlewares(authMiddleware).
    Build()

httpServer.Route(route)
```

### Route Groups

```go
userRoutes := []*server.Route{
    server.NewRouteBuilder().GET("/").Handler(getUsersHandler).Build(),
    server.NewRouteBuilder().POST("/").Handler(createUserHandler).Build(),
    server.NewRouteBuilder().GET("/{id}").Handler(getUserByIDHandler).Build(),
}

httpServer.Group("/users", userRoutes)
```

## Middleware

Middleware can be applied to the entire server or to specific routes:

```go
// Create authentication middleware
auth := middlewares.NewAuthorization(logger, tokenManager)

// Apply to the entire server
httpServer.Middleware(&server.Middleware{
    middlewares: []func(http.Handler) http.Handler{
        auth.Handle,
    },
})

// Apply to a specific route
route := server.NewRouteBuilder().
    GET("/protected").
    Handler(protectedHandler).
    Middlewares(auth.Handle).
    Build()
```

## Request Validation

Request bodies can be validated using the validator package:

```go
type UserRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var req UserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        viewmodels.NewResponseBuilder().
            Writer(w).
            BadRequest().
            Message("Invalid request body").
            Build()
        return
    }

    validator := validator.NewBodyValidator(logger)
    if httpErr := validator.Validate(req); httpErr != nil {
        viewmodels.NewResponseBuilder().
            Writer(w).
            BadRequest().
            Message(httpErr.Message).
            Details(httpErr.Details).
            Build()
        return
    }

    // Process the valid request...
}
```

## Response Building

The `viewmodels` package provides utilities for building standardized responses:

```go
// Success response
viewmodels.NewResponseBuilder().
    Writer(w).
    Ok().
    JSON(someData).
    Build()

// Error response
viewmodels.NewResponseBuilder().
    Writer(w).
    BadRequest().
    Message("Invalid input").
    Details(someDetails).
    Build()
```

## Error Handling

The package includes predefined error types and helper functions:

```go
// Create standard errors
notFoundErr := viewmodels.HTTPError{
    StatusCode: http.StatusNotFound,
    Message:    "User not found",
}

// Use helper functions
badReqErr := viewmodels.BadRequest(map[string]string{
    "email": "Invalid email format",
})

// Write to response
w.WriteHeader(badReqErr.StatusCode)
w.Write(badReqErr.ToBuffer())
```

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
