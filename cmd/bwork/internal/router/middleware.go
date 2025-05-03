package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

var globalMiddlewares []Middleware

func UseGlobalMiddleware(mw Middleware) {
	globalMiddlewares = append(globalMiddlewares, mw)
}

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

func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("📥 Ruta accedida:", r.Method, r.URL.Path)
		next(w, r)
	}
}

func EnforceJSONMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Solo se permite Content-Type: application/json", http.StatusUnsupportedMediaType)
			return
		}
		next(w, r)
	}
}

func LimitBodySizeMiddleware(maxBytes int64) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next(w, r)
		}
	}
}

type key string

const UserIDKey key = "user_id"

func InjectUserIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), UserIDKey, "12345")
		next(w, r.WithContext(ctx))
	}
}
