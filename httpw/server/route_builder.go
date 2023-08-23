package server

import "net/http"

type (
	RouteBuilder interface {
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
)

func NewRouteBuilder() RouteBuilder {
	return &routeBuilder{&Route{}}
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
