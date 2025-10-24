package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ------------------------- entry -------------------------

func runGenerate(entityType, name, moduleDir string) {
	switch entityType {
	case "controller":
		generateController(name, moduleDir)
	case "model":
		generateModel(name, moduleDir)
	case "view":
		generateView(name, moduleDir)
	case "module":
		generateModel(name, moduleDir)
		generateController(name, moduleDir)
		generateView(name, moduleDir)
	default:
		fmt.Printf("Tipo '%s' no soportado para generar\n", entityType)
	}
}

// ------------------------- helpers go.mod -------------------------

// findGoMod busca go.mod subiendo desde el cwd hasta la raíz.
func findGoMod() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		gmp := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(gmp); err == nil {
			return gmp, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("no se encontró go.mod en esta ruta o superiores")
}

// readModulePath lee la línea `module ...` del go.mod
func readModulePath() (string, error) {
	goModPath, err := findGoMod()
	if err != nil {
		return "", err
	}
	f, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	if err := sc.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("no se pudo determinar el module path desde go.mod")
}

// ------------------------- generators -------------------------

func generateController(name, moduleDir string) {
	formatted := strings.ToLower(name)
	structName := toTitle(formatted)

	content := fmt.Sprintf(`package controllers

func %sControllerLogic() string {
    return "Lógica del negocio para %s ejecutada 🚀"
}
`, structName, structName)

	dir := filepath.Join(moduleDir, "controllers")
	_ = os.MkdirAll(dir, 0o755)
	fileName := filepath.Join(dir, formatted+"_controller.go")
	_ = os.WriteFile(fileName, []byte(content), 0o644)
	fmt.Printf("Controlador '%s' generado ✅\n", fileName)
}

func generateModel(name, moduleDir string) {
	formatted := strings.ToLower(name)
	structName := toTitle(formatted)

	content := fmt.Sprintf(`package models

type %s struct {
    ID   int
    Name string
}
`, structName)

	dir := filepath.Join(moduleDir, "models")
	_ = os.MkdirAll(dir, 0o755)
	fileName := filepath.Join(dir, formatted+".go")
	_ = os.WriteFile(fileName, []byte(content), 0o644)
	fmt.Printf("Modelo '%s' generado ✅\n", fileName)
}

func generateView(name, moduleDir string) {
	formatted := strings.ToLower(name)
	handlerName := toTitle(formatted) + "View"

	content := fmt.Sprintf(`package views

import (
    "fmt"
    "net/http"
)

func %s(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Vista '%s' funcionando como endpoint 🚀")
}
`, handlerName, formatted)

	dir := filepath.Join(moduleDir, "views")
	_ = os.MkdirAll(dir, 0o755)
	fileName := filepath.Join(dir, formatted+"_view.go")
	_ = os.WriteFile(fileName, []byte(content), 0o644)
	fmt.Printf("Vista '%s' generada ✅\n", fileName)

	registerRoute(formatted, handlerName, moduleDir)
}

// ------------------------- routes.go wiring -------------------------

func registerRoute(pathName, handler, moduleDir string) {
	// modulePath es el nombre de módulo tomado de go.mod (para imports)
	modulePath, err := readModulePath()
	if err != nil {
		fmt.Println("✖ No se pudo leer go.mod:", err)
	}

	// import para views: usar slashes, no filepath separators
	viewsImport := fmt.Sprintf("%q", path.Join(modulePath, "views"))
	routesFile := filepath.Join(moduleDir, "routes.go")

	// crear si no existe
	if _, err := os.Stat(routesFile); os.IsNotExist(err) {
		base := fmt.Sprintf(`package main

import (
	"net/http"
	%s
)

func SetupRoutes(mux *http.ServeMux) {
	// Las rutas se insertan automáticamente aquí 🚀
}
`, viewsImport)
		_ = os.WriteFile(routesFile, []byte(base), 0o644)
		// seguimos para insertar la ruta
	}

	data, _ := os.ReadFile(routesFile)
	lines := strings.Split(string(data), "\n")

	// asegurar import de views
	hasImport := false
	for _, line := range lines {
		if strings.Contains(line, viewsImport) {
			hasImport = true
			break
		}
	}
	if !hasImport {
		for i, line := range lines {
			if strings.TrimSpace(line) == "import (" {
				// insertamos la línea justo después de 'import ('
				newLines := make([]string, 0, len(lines)+1)
				newLines = append(newLines, lines[:i+1]...)
				newLines = append(newLines, "\t"+viewsImport)
				newLines = append(newLines, lines[i+1:]...)
				lines = newLines
				break
			}
		}
	}

	routeLine := fmt.Sprintf("\tmux.HandleFunc(\"/%s\", views.%s)", pathName, handler)

	// evitar duplicado
	for _, line := range lines {
		if strings.Contains(line, routeLine) {
			writeRoutes(routesFile, lines)
			fmt.Printf("🔗 Ruta '/%s' ya estaba registrada en routes.go ✅\n", pathName)
			return
		}
	}

	// insertar dentro de SetupRoutes {...}
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "func SetupRoutes(") {
			// avanzar hasta la llave de apertura
			for j := i; j < len(lines); j++ {
				if strings.Contains(lines[j], "{") {
					newLines := make([]string, 0, len(lines)+1)
					newLines = append(newLines, lines[:j+1]...)
					newLines = append(newLines, routeLine)
					newLines = append(newLines, lines[j+1:]...)
					lines = newLines
					writeRoutes(routesFile, lines)
					fmt.Printf("🔗 Ruta '/%s' registrada correctamente en routes.go ✅\n", pathName)
					return
				}
			}
		}
	}

	// fallback: si no encontramos SetupRoutes, agregamos una versión mínima
	min := fmt.Sprintf(`

func SetupRoutes(mux *http.ServeMux) {
%s
}
`, routeLine)
	lines = append(lines, min)
	writeRoutes(routesFile, lines)
	fmt.Printf("🔗 Ruta '/%s' registrada correctamente en routes.go ✅\n", pathName)
}

func writeRoutes(p string, lines []string) {
	_ = os.WriteFile(p, []byte(strings.Join(lines, "\n")), 0o644)
}

// ------------------------- misc -------------------------

// strings.Title está deprecated; implementamos una simple capitalización
func toTitle(s string) string {
	if s == "" { return s }
	r := []rune(s)
	r[0] = toUpperRune(r[0])
	return string(r)
}
func toUpperRune(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - ('a' - 'A')
	}
	return r
}
