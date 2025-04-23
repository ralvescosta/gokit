# TinyHTTP Package

A lightweight HTTP server implementation for Go applications with a fluent API, built on top of the chi router.

## Overview

The `tiny_http` package provides a simple yet powerful HTTP server implementation designed for microservices and small web applications. It offers sensible defaults while allowing for easy customization through a fluent interface.

## Features

- **Lightweight**: Minimal footprint with only essential dependencies
- **Fluent API**: Easy-to-use chainable methods for configuration
- **Built-in Middleware**: Comes pre-configured with common middleware:
  - Request ID generation
  - Real IP detection
  - Panic recovery
  - Default logging
  - Content type validation
  - Response compression
  - Health check endpoint (`/health`)
- **Prometheus Integration**: Easy metrics exposure
- **Graceful Shutdown**: Proper handling of shutdown signals

## Installation

```bash
go get github.com/ralvescosta/gokit/tiny_http
```

## Usage

### Basic Example

```go
package main

import (
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/ralvescosta/gokit/configs"
    "github.com/ralvescosta/gokit/tiny_http"
)

func main() {
    // Create configuration
    cfgs := configs.NewConfigs()

    // Set up signal handling for graceful shutdown
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

    // Create and configure the server
    server := tiny_http.NewTinyServer(cfgs).
        Sig(sig).
        Route(http.MethodGet, "/api/hello", helloHandler)

    // Start the server
    if err := server.Run(); err != nil {
        panic(err)
    }
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message":"Hello, World!"}`))
}
```

### With Prometheus Metrics

```go
server := tiny_http.NewTinyServer(cfgs).
    Sig(sig).
    Prometheus().
    Route(http.MethodGet, "/api/hello", helloHandler)
```

### Adding Custom Middleware

```go
server := tiny_http.NewTinyServer(cfgs).
    Sig(sig).
    Middleware(yourCustomMiddleware).
    Route(http.MethodGet, "/api/hello", helloHandler)
```

## API Reference

### `NewTinyServer(cfgs *configs.Configs) TinyServer`

Creates a new TinyServer instance with default middleware already configured.

### TinyServer Interface

#### `Sig(sig chan os.Signal) TinyServer`

Sets a signal channel for graceful shutdown. The server will listen for signals on this channel and initiate a graceful shutdown when a signal is received.

#### `Prometheus() TinyServer`

Enables the `/metrics` endpoint for Prometheus metrics.

#### `Route(method string, path string, handler http.HandlerFunc) TinyServer`

Registers a new route with the specified HTTP method, path, and handler.

#### `Middleware(middlewares ...func(http.Handler) http.Handler) TinyServer`

Adds one or more middleware functions to the middleware stack.

#### `Run() error`

Starts the HTTP server and blocks until the server is shutdown.

## Default Configuration

- **Read Timeout**: 5 seconds
- **Write Timeout**: 10 seconds
- **Idle Timeout**: 30 seconds
- **Graceful Shutdown Timeout**: 30 seconds

## License

MIT License
