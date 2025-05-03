package router

import "net/http"

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []Middleware
}

var routes []Route

func register(method, path string, handler http.HandlerFunc, mws ...Middleware) Route {
	route := Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Middlewares: mws,
	}
	routes = append(routes, route)
	return route
}

func GET(path string, handler http.HandlerFunc, mws ...Middleware) Route {
	return register("GET", path, handler, mws...)
}
func POST(path string, handler http.HandlerFunc, mws ...Middleware) Route {
	return register("POST", path, handler, mws...)
}
func PUT(path string, handler http.HandlerFunc, mws ...Middleware) Route {
	return register("PUT", path, handler, mws...)
}
func DELETE(path string, handler http.HandlerFunc, mws ...Middleware) Route {
	return register("DELETE", path, handler, mws...)
}
