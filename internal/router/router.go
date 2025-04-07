package router

import "net/http"

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

var routes []Route

func GET(path string, handler http.HandlerFunc) Route {
	route := Route{"GET", path, handler}
	routes = append(routes, route)
	return route
}

func POST(path string, handler http.HandlerFunc) Route {
	route := Route{"POST", path, handler}
	routes = append(routes, route)
	return route
}

func PUT(path string, handler http.HandlerFunc) Route {
	route := Route{"PUT", path, handler}
	routes = append(routes, route)
	return route
}

func DELETE(path string, handler http.HandlerFunc) Route {
	route := Route{"DELETE", path, handler}
	routes = append(routes, route)
	return route
}

func ApplyRoutes(mux *http.ServeMux) {
	for _, r := range routes {
		if r.Method == "GET" {
			mux.HandleFunc(r.Path, r.Handler)
		}

	}
}
