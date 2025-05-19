package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func runGenerate(entityType, name, moduleName string) {
	switch entityType {
	case "controller":
		generateController(name, moduleName)
	case "model":
		generateModel(name, moduleName)
	case "view":
		generateView(name, moduleName)
	case "module":
		generateModel(name, moduleName)
		generateController(name, moduleName)
		generateView(name, moduleName)
	default:
		fmt.Printf("Tipo '%s' no soportado para generar\n", entityType)
	}
}

func generateController(name, moduleName string) {
	formattedName := strings.ToLower(name)
	structName := strings.Title(formattedName)

	content := fmt.Sprintf(`package controllers

func %sControllerLogic() string {
    return "Lógica del negocio para %s ejecutada 🚀"
}
`, structName, structName)

	dir := filepath.Join(moduleName, "controllers")
	os.MkdirAll(dir, 0755)
	fileName := filepath.Join(dir, formattedName+"_controller.go")
	os.WriteFile(fileName, []byte(content), 0644)
	fmt.Printf("Controlador '%s' generado ✅\n", fileName)
}

func generateModel(name, moduleName string) {
	formattedName := strings.ToLower(name)
	structName := strings.Title(formattedName)

	content := fmt.Sprintf(`package models

type %s struct {
    ID   int
    Name string
}
`, structName)

	dir := filepath.Join(moduleName, "models")
	os.MkdirAll(dir, 0755)
	fileName := filepath.Join(dir, formattedName+".go")
	os.WriteFile(fileName, []byte(content), 0644)
	fmt.Printf("Modelo '%s' generado ✅\n", fileName)
}

func generateView(name, moduleName string) {
	formattedName := strings.ToLower(name)
	handlerName := strings.Title(formattedName) + "View"

	content := fmt.Sprintf(`package views

import (
    "fmt"
    "net/http"
)

func %s(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Vista '%s' funcionando como endpoint 🚀")
}
`, handlerName, formattedName)

	dir := filepath.Join(moduleName, "views")
	os.MkdirAll(dir, 0755)
	fileName := filepath.Join(dir, formattedName+"_view.go")
	os.WriteFile(fileName, []byte(content), 0644)
	fmt.Printf("Vista '%s' generada ✅\n", fileName)

	registerRoute(formattedName, handlerName, moduleName)
}

func registerRoute(path, handler, moduleName string) {
	routesFile := filepath.Join(moduleName, moduleName, "routes.go") // ✅ corregido

	viewImport := fmt.Sprintf("\"%s/views\"", moduleName) // ✅ ejemplo: "back/views"

	// Si no existe routes.go, crearlo con base mínima
	if _, err := os.Stat(routesFile); os.IsNotExist(err) {
		base := fmt.Sprintf(`package main

import (
	"net/http"
	%s
)

func SetupRoutes(mux *http.ServeMux) {
	// Aquí se registrarán las rutas automáticamente 🚀
}
`, viewImport)
		os.WriteFile(routesFile, []byte(base), 0644)
		return
	}

	// Leer archivo existente
	data, _ := os.ReadFile(routesFile)
	lines := strings.Split(string(data), "\n")

	// Agregar import si falta
	hasImport := false
	for _, line := range lines {
		if strings.Contains(line, viewImport) {
			hasImport = true
			break
		}
	}
	if !hasImport {
		for i, line := range lines {
			if strings.Contains(line, "import (") {
				lines = append(lines[:i+1], append([]string{viewImport}, lines[i+1:]...)...)
				break
			}
		}
	}

	// Agregar ruta si falta
	routeLine := fmt.Sprintf("\tmux.HandleFunc(\"/%s\", views.%s)", path, handler)
	hasRoute := false
	for _, line := range lines {
		if strings.Contains(line, routeLine) {
			hasRoute = true
			break
		}
	}
	if !hasRoute {
		for i, line := range lines {
			if strings.Contains(line, "SetupRoutes") && strings.Contains(line, "{") {
				lines = append(lines[:i+1], append([]string{routeLine}, lines[i+1:]...)...)
				break
			}
		}
	}

	// Guardar archivo actualizado
	output := strings.Join(lines, "\n")
	os.WriteFile(routesFile, []byte(output), 0644)

	fmt.Printf("🔗 Ruta '/%s' registrada correctamente en routes.go ✅\n", path)
}
