// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package server

import "net/http"

type (
	RouteBuilder interface {
		GET(path string) RouteBuilder
		POST(path string) RouteBuilder
		PUT(path string) RouteBuilder
		PATCH(path string) RouteBuilder
		DELETE(path string) RouteBuilder
		Path(p string) RouteBuilder
		Method(m string) RouteBuilder
		Handler(h http.HandlerFunc) RouteBuilder
		Middlewares(m ...func(http.Handler) http.Handler) RouteBuilder
		Build() *Route
	}

	routeBuilder struct {
		r *Route
	}

	Route struct {
		method      string
		path        string
		handler     http.HandlerFunc
		middlewares []func(http.Handler) http.Handler
	}

	Middleware struct {
		middlewares []func(http.Handler) http.Handler
	}
)

func NewRouteBuilder() RouteBuilder {
	return &routeBuilder{&Route{}}
}

func (b *routeBuilder) GET(path string) RouteBuilder {
	b.r.path = path
	b.r.method = http.MethodGet
	return b
}

func (b *routeBuilder) POST(path string) RouteBuilder {
	b.r.path = path
	b.r.method = http.MethodPost
	return b
}

func (b *routeBuilder) PUT(path string) RouteBuilder {
	b.r.path = path
	b.r.method = http.MethodPut
	return b
}

func (b *routeBuilder) PATCH(path string) RouteBuilder {
	b.r.path = path
	b.r.method = http.MethodPatch
	return b
}

func (b *routeBuilder) DELETE(path string) RouteBuilder {
	b.r.path = path
	b.r.method = http.MethodDelete
	return b
}

func (b *routeBuilder) Path(p string) RouteBuilder {
	b.r.path = p
	return b
}

func (b *routeBuilder) Method(m string) RouteBuilder {
	b.r.method = m
	return b
}

func (b *routeBuilder) Handler(h http.HandlerFunc) RouteBuilder {
	b.r.handler = h
	return b
}

func (b *routeBuilder) Middlewares(m ...func(http.Handler) http.Handler) RouteBuilder {
	b.r.middlewares = m
	return b
}

func (b *routeBuilder) Build() *Route {
	return b.r
}
