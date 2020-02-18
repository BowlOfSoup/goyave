package goyave

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/System-Glitch/goyave/v2/validation"
)

type routerDefinition struct {
	prefix     string // Empty for main router
	middleware []Middleware
	routes     []*routeDefinition
	subrouters []*routerDefinition
}

type routeDefinition struct {
	uri        string
	methods    string
	name       string
	handler    Handler
	rules      validation.RuleSet
	middleware []Middleware
}

var handler Handler = func(response *Response, request *Request) {
	response.Status(200)
}

var sampleRouteDefinition *routerDefinition = &routerDefinition{
	prefix:     "",
	middleware: []Middleware{},
	routes: []*routeDefinition{
		{
			uri:     "/hello",
			methods: "GET",
			name:    "hello",
			handler: handler,
			rules:   nil,
		},
		{
			uri:     "/{param}",
			methods: "POST",
			name:    "param",
			handler: handler,
			rules:   nil,
		},
	},
	subrouters: []*routerDefinition{
		{
			prefix:     "/product",
			middleware: []Middleware{},
			routes: []*routeDefinition{
				{
					uri:     "/",
					methods: "GET",
					name:    "product.index",
					handler: handler,
					rules:   nil,
				},
				{
					uri:     "/",
					methods: "POST",
					name:    "product.store",
					handler: handler,
					rules:   nil,
				},
				{
					uri:     "/{id:[0-9]+}",
					methods: "GET",
					name:    "product.show",
					handler: handler,
					rules:   nil,
				},
				{
					uri:     "/{id:[0-9]+}",
					methods: "PUT|PATCH",
					name:    "product.update",
					handler: handler,
					rules:   nil,
				},
				{
					uri:     "/{id:[0-9]+}",
					methods: "DELETE",
					name:    "product.destroy",
					handler: handler,
					rules:   nil,
				},
			},
			subrouters: []*routerDefinition{},
		},
	},
}

var sampleRequests []*http.Request = []*http.Request{
	httptest.NewRequest("GET", "/", nil), // 404
	httptest.NewRequest("GET", "/hello", nil),
	httptest.NewRequest("POST", "/world", nil),
	httptest.NewRequest("GET", "/product", nil),
	httptest.NewRequest("POST", "/product", nil),     // TODO body and validation
	httptest.NewRequest("GET", "/product/test", nil), // 404
	httptest.NewRequest("GET", "/product/1", nil),
	httptest.NewRequest("PUT", "/product/1", nil),
	httptest.NewRequest("DELETE", "/product/1", nil),
}

func registerAll(def *routerDefinition) *Router {
	main := newRouter()
	registerRouter(main, def)
	return main
}

func registerRouter(router *Router, def *routerDefinition) {
	for _, subdef := range def.subrouters {
		subrouter := router.Subrouter(subdef.prefix)
		registerRouter(subrouter, subdef)
	}
	for _, routeDef := range def.routes {
		router.registerRoute(routeDef.methods, routeDef.uri, routeDef.handler, routeDef.rules).Name(routeDef.name)
	}
}

func BenchmarkRouteRegistration(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		registerAll(sampleRouteDefinition)
		regexCache = nil
	}
}

func BenchmarkRootLevelNotFound(b *testing.B) {
	router := setupRouteBench(b)

	for n := 0; n < b.N; n++ {
		router.match(sampleRequests[0], &routeMatch{})
	}
}

func BenchmarkRootLevelMatch(b *testing.B) {
	router := setupRouteBench(b)

	for n := 0; n < b.N; n++ {
		router.match(sampleRequests[1], &routeMatch{})
	}
}

func BenchmarkRootLevelPostMatch(b *testing.B) {
	router := setupRouteBench(b)

	for n := 0; n < b.N; n++ {
		router.match(sampleRequests[2], &routeMatch{})
	}
}

func BenchmarkSubrouterMatch(b *testing.B) {
	router := setupRouteBench(b)

	for n := 0; n < b.N; n++ {
		router.match(sampleRequests[3], &routeMatch{})
	}
}

func BenchmarkSubrouterPostMatch(b *testing.B) {
	router := setupRouteBench(b)

	for n := 0; n < b.N; n++ {
		router.match(sampleRequests[4], &routeMatch{})
	}
}

func BenchmarkSubrouterNotFound(b *testing.B) {
	router := setupRouteBench(b)

	for n := 0; n < b.N; n++ {
		router.match(sampleRequests[5], &routeMatch{})
	}
}

func BenchmarkParamMatch(b *testing.B) {
	router := setupRouteBench(b)

	for n := 0; n < b.N; n++ {
		router.match(sampleRequests[6], &routeMatch{})
	}
}

func BenchmarkParamPutMatch(b *testing.B) {
	router := setupRouteBench(b)

	for n := 0; n < b.N; n++ {
		router.match(sampleRequests[7], &routeMatch{})
	}
}

func BenchmarkParamDeleteMatch(b *testing.B) {
	router := setupRouteBench(b)

	for n := 0; n < b.N; n++ {
		router.match(sampleRequests[8], &routeMatch{})
	}
}

func setupRouteBench(b *testing.B) *Router {
	router := registerAll(sampleRouteDefinition)
	regexCache = nil
	b.ReportAllocs()
	runtime.GC()
	defer b.ResetTimer()
	return router
}
