package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed internal/router/router.go
var routerSource []byte

func runInit() {
	fmt.Println("Inicializando proyecto BWORK...")

	// Crear estructura de carpetas
	os.Mkdir("bwork_modules", 0755)
	os.MkdirAll("app/controllers", 0755)
	os.MkdirAll("app/models", 0755)
	os.MkdirAll("app/views", 0755)
	os.MkdirAll("app/config", 0755)

	// Archivos de configuración
	os.WriteFile("bwork.json", []byte("{\n  \"name\": \"mi-backend\",\n  \"version\": \"0.0.1\"\n}"), 0644)
	os.WriteFile(".env", []byte("DB_HOST=localhost\nDB_PORT=3306\nDB_USER=root\nDB_PASSWORD=password\nDB_NAME=bworkdb\n"), 0644)
	os.WriteFile(".gitignore", []byte(".env\nbwork_modules/\n*.log\n*.tmp\n*.out\n"), 0644)

	// Ejecutar: go mod init app
	cmd := exec.Command("go", "mod", "init", "app")
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

	os.Mkdir("app", 0755)
	os.WriteFile(filepath.Join("app", "main.go"), []byte(mainContent), 0644)
	os.WriteFile(filepath.Join("app", "routes.go"), []byte(routesContent), 0644)

	// Copiar módulo router a bwork_modules/router
	copyRouterModule()

	// Crear README.md
	readmeContent := "# 🚀 Proyecto creado con BWORK\n\n" +
		"Este backend fue generado con [BWORK](https://github.com/alun-dra/bwork), un framework CLI para construir APIs Go de forma rápida y modular.\n\n" +
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
		"Este proyecto fue generado con 💡 usando [BWORK CLI](https://github.com/alun-dra/bwork).\n" +
		"Si quieres colaborar o sugerir mejoras, ¡haz un PR o crea un issue!\n"

	os.WriteFile("README.md", []byte(readmeContent), 0644)

	fmt.Println("✅ Proyecto BWORK inicializado con éxito.")
}

func copyRouterModule() {
	destDir := filepath.Join("bwork_modules", "router")
	dest := filepath.Join(destDir, "router.go")

	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		fmt.Println("❌ No se pudo crear el directorio del módulo router:", err)
		return
	}

	// Archivos fuente del router
	files := []string{
		"internal/router/core.go",
		"internal/router/middleware.go",
		"internal/router/utils.go",
		"internal/router/router.go", // este archivo debe tener solo ApplyRoutes
	}

	var fullRouterCode []byte
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("❌ Error leyendo %s: %v\n", file, err)
			return
		}
		fullRouterCode = append(fullRouterCode, append(content, []byte("\n\n")...)...)
	}

	err = os.WriteFile(dest, fullRouterCode, 0644)
	if err != nil {
		fmt.Println("❌ Error al escribir el archivo router.go combinado:", err)
		return
	}

	fmt.Println("📦 Módulo 'router' copiado y ensamblado en .bwork_modules/router ✅")
}
