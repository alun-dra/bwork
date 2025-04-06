package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func runInit() {
	fmt.Println("Inicializando proyecto BWORK...")

	// Crear estructura de carpetas
	os.Mkdir(".bwork_modules", 0755)
	os.MkdirAll("app/controllers", 0755)
	os.MkdirAll("app/models", 0755)
	os.MkdirAll("app/views", 0755)
	os.MkdirAll("app/config", 0755)

	// Archivos de configuración
	os.WriteFile("bwork.json", []byte("{\n  \"name\": \"mi-backend\",\n  \"version\": \"0.0.1\"\n}"), 0644)
	os.WriteFile(".env", []byte("DB_HOST=localhost\nDB_PORT=3306\nDB_USER=root\nDB_PASSWORD=password\nDB_NAME=bworkdb\n"), 0644)
	os.WriteFile(".gitignore", []byte(".env\n.bwork_modules/\n*.log\n*.tmp\n*.out\n"), 0644)

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

	// main.go con servidor funcional
	mainContent := `package main

import (
	"fmt"
	"log"
	"net/http"
	"app"
)

func main() {
	fmt.Println("🚀 Servidor iniciado en http://localhost:8081")

	mux := http.NewServeMux()
	app.SetupRoutes(mux)

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatal("❌ Error al iniciar el servidor:", err)
	}
}
`
	os.Mkdir("app", 0755)
	os.WriteFile(filepath.Join("app", "main.go"), []byte(mainContent), 0644)

	fmt.Println("Proyecto inicializado ✅")
}
