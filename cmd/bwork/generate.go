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

	// Agregar automáticamente la ruta
	registerRoute(formattedName, handlerName)
}

func registerRoute(path, handler string) {
	routesFile := "app/routes.go"

	// Si no existe, crearlo con estructura base
	if _, err := os.Stat(routesFile); os.IsNotExist(err) {
		base := `package main

import (
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
}
`
		os.WriteFile(routesFile, []byte(base), 0644)
	}

	// Leer archivo
	data, _ := os.ReadFile(routesFile)
	text := string(data)

	// Agregar import si no existe
	if !strings.Contains(text, `"app/views"`) {
		text = strings.Replace(text, "import (", "import (\n\t\"app/views\"", 1)
	}

	// Preparar línea de ruta
	routeLine := fmt.Sprintf("mux.HandleFunc(\"/%s\", views.%s)", path, handler)

	// Insertar dentro de SetupRoutes si no está ya
	if !strings.Contains(text, routeLine) {
		text = strings.Replace(text, "{", "{\n\t"+routeLine, 1)
	}

	// Guardar archivo actualizado
	os.WriteFile(routesFile, []byte(text), 0644)
	fmt.Printf("🔗 Ruta '/%s' registrada en app/routes.go\n", path)
}
