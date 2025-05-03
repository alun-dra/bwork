package router

import "net/http"

func applyMiddleware(handler http.HandlerFunc, mws []Middleware) http.HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}

func Chain(handler http.HandlerFunc, mws ...Middleware) http.HandlerFunc {
	return applyMiddleware(handler, mws)
}
