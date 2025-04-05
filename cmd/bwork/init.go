package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func runInit() {
    fmt.Println("Inicializando proyecto BWORK...")

    os.Mkdir(".bwork_modules", 0755)
    os.MkdirAll("app/controllers", 0755)
    os.MkdirAll("app/models", 0755)
    os.MkdirAll("app/views", 0755)
    os.MkdirAll("app/config", 0755)

    os.WriteFile("bwork.json", []byte("{\n  \"name\": \"mi-backend\",\n  \"version\": \"0.0.1\"\n}"), 0644)
    os.WriteFile(".env", []byte("DB_HOST=localhost\nDB_PORT=3306\nDB_USER=root\nDB_PASSWORD=password\nDB_NAME=bworkdb\n"), 0644)
    os.WriteFile(".gitignore", []byte(".env\n.bwork_modules/\n*.log\n*.tmp\n*.out\n"), 0644)

    mainContent := `package main

import "fmt"

func main() {
    fmt.Println("Bienvenido a tu backend con BWORK 🚀")
}
`
    os.Mkdir("app", 0755)
    os.WriteFile(filepath.Join("app", "main.go"), []byte(mainContent), 0644)

    fmt.Println("Proyecto inicializado ✅")
}
