package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func runInstall(lib string) {
	fmt.Printf("🔧 Instalando librería '%s'...\n", lib)

	var content string

	switch lib {
	case "router":
		content = `package router

import "fmt"

func SetupRouter() {
    fmt.Println("Router configurado 🚀")
}
`
	case "controller":
		content = `package controller

import "fmt"

func DefaultController() {
    fmt.Println("Controlador por defecto listo ✅")
}
`
	case "orm":
		content = `package orm

import "fmt"

func Connect() {
    fmt.Println("Conexión al ORM establecida 🔗")
}
`
	default:
		fmt.Println("❌ Librería no reconocida. Usa: router, controller, orm")
		return
	}

	dir := filepath.Join("bwork_modules", lib)
	os.MkdirAll(dir, 0755)

	file := filepath.Join(dir, lib+".go")
	os.WriteFile(file, []byte(content), 0644)

	fmt.Println("✅ Librería instalada en", file)
}
