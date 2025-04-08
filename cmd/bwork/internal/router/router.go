package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []Middleware
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

var routes []Route
var globalMiddlewares []Middleware

// ----------------- REGISTRO DE RUTAS -----------------

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

// ----------------- MIDDLEWARE -----------------

func UseGlobalMiddleware(mw Middleware) {
	globalMiddlewares = append(globalMiddlewares, mw)
}

func applyMiddleware(handler http.HandlerFunc, mws []Middleware) http.HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}
	return handler
}

func Chain(handler http.HandlerFunc, mws ...Middleware) http.HandlerFunc {
	return applyMiddleware(handler, mws)
}

// ----------------- REGISTRO DE TODAS LAS RUTAS -----------------

func ApplyRoutes(mux *http.ServeMux) {
	for _, r := range routes {
		if r.Method == "GET" || r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
			handler := applyMiddleware(r.Handler, append(globalMiddlewares, r.Middlewares...))
			mux.HandleFunc(r.Path, handler)
		}
	}
}

// ----------------- MIDDLEWARES EMBEBIDOS -----------------

// AuthMiddleware verifica token de autorización simple
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer mi-secreto" {
			http.Error(w, "No autorizado", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// LogMiddleware imprime la ruta accedida
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("📥 Ruta accedida:", r.Method, r.URL.Path)
		next(w, r)
	}
}

// EnforceJSONMiddleware rechaza peticiones sin content-type JSON
func EnforceJSONMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Solo se permite Content-Type: application/json", http.StatusUnsupportedMediaType)
			return
		}
		next(w, r)
	}
}

// LimitBodySizeMiddleware limita el tamaño del body
func LimitBodySizeMiddleware(maxBytes int64) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next(w, r)
		}
	}
}

// InjectUserIDMiddleware inserta un userID falso al contexto
type key string

const UserIDKey key = "user_id"

func InjectUserIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), UserIDKey, "12345")
		next(w, r.WithContext(ctx))
	}
}
