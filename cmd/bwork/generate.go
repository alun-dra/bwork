package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func runGenerate(entityType, name string) {
	switch entityType {
	case "controller":
		generateController(name)
	case "model":
		generateModel(name)
	case "view":
		generateView(name)
	case "module":
		generateModel(name)
		generateController(name)
		generateView(name)
	default:
		fmt.Printf("Tipo '%s' no soportado para generar\n", entityType)
	}
}

func generateController(name string) {
	formattedName := strings.ToLower(name)
	structName := strings.Title(formattedName)

	content := fmt.Sprintf(`package controllers

func %sControllerLogic() string {
    return "Lógica del negocio para %s ejecutada 🚀"
}
`, structName, structName)

	dir := "app/controllers"
	os.MkdirAll(dir, 0755)
	fileName := filepath.Join(dir, formattedName+"_controller.go")
	os.WriteFile(fileName, []byte(content), 0644)
	fmt.Printf("Controlador '%s' generado ✅\n", fileName)
}

func generateModel(name string) {
	formattedName := strings.ToLower(name)
	structName := strings.Title(formattedName)

	content := fmt.Sprintf(`package models

type %s struct {
    ID   int
    Name string
}
`, structName)

	dir := "app/models"
	os.MkdirAll(dir, 0755)
	fileName := filepath.Join(dir, formattedName+".go")
	os.WriteFile(fileName, []byte(content), 0644)
	fmt.Printf("Modelo '%s' generado ✅\n", fileName)
}

func generateView(name string) {
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

	dir := "app/views"
	os.MkdirAll(dir, 0755)
	fileName := filepath.Join(dir, formattedName+"_view.go")
	os.WriteFile(fileName, []byte(content), 0644)
	fmt.Printf("Vista '%s' generada ✅\n", fileName)

	registerRoute(formattedName, handlerName)
}

func registerRoute(path, handler string) {
	routesFile := "app/routes.go"

	// Detectar el módulo desde go.mod
	moduleName := "app"
	if data, err := os.ReadFile("go.mod"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "module ") {
				moduleName = strings.TrimSpace(strings.TrimPrefix(line, "module"))
				break
			}
		}
	}

	viewImport := fmt.Sprintf("\"%s/views\"", moduleName)

	// Si no existe routes.go, crearlo con bloque de import
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
		fmt.Println("✅ Archivo routes.go creado")
		return
	}

	// Leer contenido
	data, _ := os.ReadFile(routesFile)
	lines := strings.Split(string(data), "\n")

	// 1. Verificar si ya tiene el import
	hasImport := false
	for _, line := range lines {
		if strings.Contains(line, viewImport) {
			hasImport = true
			break
		}
	}

	// 2. Asegurar que sea import (...) si aún no lo es
	for i, line := range lines {
		if strings.HasPrefix(line, "import ") && !strings.Contains(line, "(") {
			singleImport := strings.TrimSpace(strings.TrimPrefix(line, "import"))
			lines[i] = "import (\n\t" + singleImport + "\n)"
			break
		}
	}

	// 3. Agregar el import views si no existe
	if !hasImport {
		for i, line := range lines {
			if strings.HasPrefix(line, "import (") {
				lines = append(lines[:i+1], append([]string{"\t" + viewImport}, lines[i+1:]...)...)
				break
			}
		}
	}

	// 4. Agregar la línea mux.HandleFunc si no está
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

	// Guardar el archivo
	output := strings.Join(lines, "\n")
	os.WriteFile(routesFile, []byte(output), 0644)

	fmt.Printf("🔗 Ruta '/%s' registrada correctamente en routes.go ✅\n", path)
}
