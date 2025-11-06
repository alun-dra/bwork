package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed internal/router/router.go
var routerSource []byte

//go:embed internal/router/core.go
var coreSource []byte

//go:embed internal/router/middleware.go
var middlewareSource []byte

//go:embed internal/router/utils.go
var utilsSource []byte

//go:embed internal/router/router.go
var routerMainSource []byte

func runInit(moduleName string) {
	fmt.Println("Inicializando proyecto BWORK...")

	// Crear estructura de carpetas
	os.MkdirAll(filepath.Join(moduleName, "bwork_modules", "router"), 0755)
	os.MkdirAll(filepath.Join(moduleName, "controllers"), 0755)
	os.MkdirAll(filepath.Join(moduleName, "models"), 0755)
	os.MkdirAll(filepath.Join(moduleName, "views"), 0755)
	os.MkdirAll(filepath.Join(moduleName, "config"), 0755)

	// Archivos de configuración
	os.WriteFile("bwork.json", []byte("{\n  \"name\": \"mi-backend\",\n  \"version\": \"0.0.1\"\n}"), 0644)
	os.WriteFile(".env", []byte("DB_HOST=localhost\nDB_PORT=3306\nDB_USER=root\nDB_PASSWORD=password\nDB_NAME=bworkdb\n"), 0644)
	os.WriteFile(".gitignore", []byte(".env\nbwork_modules/\n*.log\n*.tmp\n*.out\n"), 0644)

	// Ejecutar: go mod init app
	cmd := exec.Command("go", "mod", "init", moduleName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Error al crear go.mod:", err)
	} else {
		fmt.Println("📦 Módulo Go inicializado como 'app'")
	}

	// Ejecutar: go mod tidy
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	err = tidyCmd.Run()
	if err != nil {
		fmt.Println("❌ Error al ejecutar 'go mod tidy':", err)
	} else {
		fmt.Println("✅ Dependencias y módulo Go configurados correctamente")
	}

	// Crear app/main.go
	mainContent := `package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("🚀 Servidor iniciado en http://localhost:8081")

	mux := http.NewServeMux()
	SetupRoutes(mux)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatal("❌ Error al iniciar el servidor:", err)
	}
}`

	// Crear app/routes.go
	routesContent := `package main

import (
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	// Aquí se registrarán las rutas automáticamente 🚀
}`

	os.Mkdir(moduleName, 0755)
	os.WriteFile(filepath.Join(moduleName, "main.go"), []byte(mainContent), 0644)
	os.WriteFile(filepath.Join(moduleName, "routes.go"), []byte(routesContent), 0644)

	copyRouterModule(moduleName)
	createRouterGoMod(moduleName)
	addReplaceDirective(moduleName)
	addRequireDirective(moduleName)
	addRouterRequireDirective(moduleName)
	generateCorsFile(moduleName)

	// Crear README.md
	readmeContent := "# 🚀 Proyecto creado con BWORK\n\n" +
		"Este backend fue generado con [BWORK](https://github.com/alun-dra/bwork), un framework para construir APIs Go de forma rápida y modular.\n\n" +
		"---\n\n" +
		"## 📁 Estructura del proyecto\n\n" +
		"```bash\n" +
		".\n" +
		"├── app/\n" +
		"│   ├── controllers/\n" +
		"│   ├── models/\n" +
		"│   ├── views/\n" +
		"│   ├── config/\n" +
		"│   ├── main.go\n" +
		"│   └── routes.go\n" +
		"├── .env\n" +
		"├── .gitignore\n" +
		"├── bwork.json\n" +
		"├── go.mod\n" +
		"```\n\n" +
		"---\n\n" +
		"## 🛠 Comandos disponibles\n\n" +
		"```bash\n" +
		"bwork init                        # Inicializa un nuevo proyecto\n" +
		"bwork generate model <name>      # Crea un modelo\n" +
		"bwork generate controller <name> # Crea un controlador\n" +
		"bwork generate view <name>       # Crea una vista (endpoint)\n" +
		"bwork generate module <name>     # Crea model + controller + view\n" +
		"bwork run                         # Ejecuta el servidor en http://localhost:8081\n" +
		"```\n\n" +
		"---\n\n" +
		"## ⚡ Ejemplo rápido\n\n" +
		"```bash\n" +
		"bwork generate module usuario\n" +
		"bwork run\n" +
		"```\n\n" +
		"Luego visita: [http://localhost:8081/usuario](http://localhost:8081/usuario)\n\n" +
		"---\n\n" +
		"## 🧪 Variables de entorno (.env)\n\n" +
		"```env\n" +
		"DB_HOST=localhost\n" +
		"DB_PORT=3306\n" +
		"DB_USER=root\n" +
		"DB_PASSWORD=password\n" +
		"DB_NAME=bworkdb\n" +
		"```\n\n" +
		"---\n\n" +
		"## 🧩 Personalización del módulo\n\n" +
		"Si deseas cambiar el nombre del módulo `app`, edita el archivo `go.mod`:\n\n" +
		"```go\n" +
		"module github.com/tunombre/miproyecto\n" +
		"```\n\n" +
		"Y actualiza los imports en `main.go`, `routes.go`, etc.\n\n" +
		"---\n\n" +
		"## 🧠 Contribuye\n\n" +
		"\n" +
		"Si quieres colaborar o sugerir mejoras, ¡haz un PR o crea un issue!\n"

	os.WriteFile("README.md", []byte(readmeContent), 0644)

	fmt.Println("✅ Proyecto BWORK inicializado con éxito.")
}

func copyRouterModule(moduleName string) {
	destDir := filepath.Join(moduleName, "bwork_modules", "router")

	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		fmt.Println("❌ No se pudo crear el directorio del módulo router:", err)
		return
	}

	files := []struct {
		Name    string
		Content []byte
	}{
		{"core.go", coreSource},
		{"middleware.go", middlewareSource},
		{"utils.go", utilsSource},
		{"router.go", routerMainSource},
	}

	for _, file := range files {
		destPath := filepath.Join(destDir, file.Name)
		err := os.WriteFile(destPath, file.Content, 0644)
		if err != nil {
			fmt.Printf("❌ Error al escribir %s: %v\n", file.Name, err)
			return
		}
	}

	fmt.Println("📦 Módulo 'router' copiado por separado en bwork_modules/router ✅")
}

func createRouterGoMod(moduleName string) {
	rootModPath := "go.mod"
	content, err := os.ReadFile(rootModPath)
	if err != nil {
		fmt.Println("❌ No se pudo leer la versión de Go desde go.mod raíz:", err)
		return
	}

	goVersion := "1.20"
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "go ") {
			goVersion = strings.TrimSpace(strings.TrimPrefix(line, "go "))
			break
		}
	}

	modContent := fmt.Sprintf("module %s/bwork_modules/router\n\ngo %s\n", moduleName, goVersion)
	goModPath := filepath.Join(moduleName, "bwork_modules", "router", "go.mod")

	err = os.WriteFile(goModPath, []byte(modContent), 0644)
	if err != nil {
		fmt.Println("❌ No se pudo crear go.mod del submódulo router:", err)
		return
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = filepath.Join(moduleName, "bwork_modules", "router")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("❌ Error al ejecutar 'go mod tidy' en submódulo router:", err)
		return
	}

	fmt.Println("📦 go.mod creado para submódulo router con tidy ✅")
}

func addReplaceDirective(moduleName string) {
	rootGoMod := "go.mod"
	replaceLine := fmt.Sprintf("\nreplace %s/bwork_modules/router => ./%s/bwork_modules/router\n", moduleName, moduleName)

	content, err := os.ReadFile(rootGoMod)
	if err != nil {
		fmt.Println("❌ No se pudo leer go.mod raíz para añadir replace:", err)
		return
	}

	if !strings.Contains(string(content), fmt.Sprintf("replace %s/bwork_modules/router", moduleName)) {
		newContent := string(content) + replaceLine
		err = os.WriteFile(rootGoMod, []byte(newContent), 0644)
		if err != nil {
			fmt.Println("❌ No se pudo escribir en go.mod raíz:", err)
			return
		}
		fmt.Println("🔁 Línea 'replace' añadida a go.mod raíz ✅")
	}
}

func addRequireDirective(moduleName string) {
	rootGoMod := "go.mod"
	requireLine := fmt.Sprintf("\nrequire %s/bwork_modules/router v0.0.0-00010101000000-000000000000\n", moduleName)

	content, err := os.ReadFile(rootGoMod)
	if err != nil {
		fmt.Println("❌ No se pudo leer go.mod raíz para añadir require:", err)
		return
	}

	if !strings.Contains(string(content), fmt.Sprintf("require %s/bwork_modules/router", moduleName)) {
		newContent := string(content) + requireLine
		err = os.WriteFile(rootGoMod, []byte(newContent), 0644)
		if err != nil {
			fmt.Println("❌ No se pudo escribir en go.mod raíz para añadir require:", err)
			return
		}
		fmt.Println("📦 Línea 'require' añadida a go.mod raíz ✅")
	}
}

func addRouterRequireDirective(moduleName string) {
	rootGoMod := "go.mod"
	requireLine := fmt.Sprintf("\nrequire %s/bwork_modules/router v0.0.0\n", moduleName)

	content, err := os.ReadFile(rootGoMod)
	if err != nil {
		fmt.Println("❌ No se pudo leer go.mod raíz para añadir require:", err)
		return
	}

	if !strings.Contains(string(content), fmt.Sprintf("require %s/bwork_modules/router", moduleName)) {
		newContent := string(content) + requireLine
		err = os.WriteFile(rootGoMod, []byte(newContent), 0644)
		if err != nil {
			fmt.Println("❌ No se pudo escribir la línea require:", err)
			return
		}
		fmt.Println("📥 Línea 'require' añadida a go.mod raíz ✅")
	}
}

func generateCorsFile(moduleName string) {
	corsContent := `package config

import "net/http"

func SetupCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
`

	configDir := filepath.Join(moduleName, "config")
	os.MkdirAll(configDir, 0755)

	filePath := filepath.Join(configDir, "cors.go")
	err := os.WriteFile(filePath, []byte(corsContent), 0644)
	if err != nil {
		fmt.Println("❌ Error al generar archivo CORS:", err)
		return
	}

	fmt.Println("🛡️ Archivo 'cors.go' generado correctamente en config/ ✅")
}
