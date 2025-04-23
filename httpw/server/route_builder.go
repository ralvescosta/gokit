// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package server provides components for building and managing HTTP servers.
package server

import "net/http"

type (
	// RouteBuilder defines the interface for building HTTP routes.
	// It follows the builder pattern to fluently construct Route objects.
	RouteBuilder interface {
		// GET creates a route with the GET HTTP method and specified path.
		GET(path string) RouteBuilder

		// POST creates a route with the POST HTTP method and specified path.
		POST(path string) RouteBuilder

		// PUT creates a route with the PUT HTTP method and specified path.
		PUT(path string) RouteBuilder

		// PATCH creates a route with the PATCH HTTP method and specified path.
		PATCH(path string) RouteBuilder

		// DELETE creates a route with the DELETE HTTP method and specified path.
		DELETE(path string) RouteBuilder

		// Path sets the URL path for the route.
		Path(p string) RouteBuilder

		// Method sets the HTTP method for the route.
		Method(m string) RouteBuilder

		// Handler sets the function that handles requests to this route.
		Handler(h http.HandlerFunc) RouteBuilder

		// Middlewares adds middleware functions to the route.
		Middlewares(m ...func(http.Handler) http.Handler) RouteBuilder

		// Build constructs and returns a Route object with the configured settings.
		Build() *Route
	}

	// routeBuilder implements the RouteBuilder interface.
	routeBuilder struct {
		r *Route
	}

	// Route represents an HTTP route with its method, path, handler, and middlewares.
	// It encapsulates all the information needed to register a route with the HTTP server.
	Route struct {
		method      string
		path        string
		handler     http.HandlerFunc
		middlewares []func(http.Handler) http.Handler
	}

	// Middleware holds a list of middleware functions that can be applied to routes.
	Middleware struct {
		middlewares []func(http.Handler) http.Handler
	}
)

// NewRouteBuilder creates a new RouteBuilder instance for fluent route construction.
func NewRouteBuilder() RouteBuilder {
	return &routeBuilder{&Route{}}
}

// GET creates a route with the GET HTTP method and specified path.
func (b *routeBuilder) GET(path string) RouteBuilder {
	b.r.path = path
	b.r.method = http.MethodGet
	return b
}

// POST creates a route with the POST HTTP method and specified path.
func (b *routeBuilder) POST(path string) RouteBuilder {
	b.r.path = path
	b.r.method = http.MethodPost
	return b
}

// PUT creates a route with the PUT HTTP method and specified path.
func (b *routeBuilder) PUT(path string) RouteBuilder {
	b.r.path = path
	b.r.method = http.MethodPut
	return b
}

// PATCH creates a route with the PATCH HTTP method and specified path.
func (b *routeBuilder) PATCH(path string) RouteBuilder {
	b.r.path = path
	b.r.method = http.MethodPatch
	return b
}

// DELETE creates a route with the DELETE HTTP method and specified path.
func (b *routeBuilder) DELETE(path string) RouteBuilder {
	b.r.path = path
	b.r.method = http.MethodDelete
	return b
}

// Path sets the URL path for the route.
func (b *routeBuilder) Path(p string) RouteBuilder {
	b.r.path = p
	return b
}

// Method sets the HTTP method for the route.
func (b *routeBuilder) Method(m string) RouteBuilder {
	b.r.method = m
	return b
}

// Handler sets the function that handles requests to this route.
func (b *routeBuilder) Handler(h http.HandlerFunc) RouteBuilder {
	b.r.handler = h
	return b
}

// Middlewares adds middleware functions to the route.
func (b *routeBuilder) Middlewares(m ...func(http.Handler) http.Handler) RouteBuilder {
	b.r.middlewares = m
	return b
}

// Build constructs and returns a Route object with the configured settings.
func (b *routeBuilder) Build() *Route {
	return b.r
}
