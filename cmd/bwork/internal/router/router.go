package router

import "net/http"

func ApplyRoutes(mux *http.ServeMux) {
	for _, r := range routes {
		if r.Method == "GET" || r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
			handler := applyMiddleware(r.Handler, append(globalMiddlewares, r.Middlewares...))
			mux.HandleFunc(r.Path, handler)
		}
	}
}
