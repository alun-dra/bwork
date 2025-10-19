package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// runServer ejecuta el servidor Go del proyecto generado por bwork init {nombre}
func runServer() {
	projectDir := detectProjectDir()
	if projectDir == "" {
		fmt.Println("❌ No se encontró un proyecto válido para ejecutar (carpeta con main.go).")
		return
	}

	mainFile := filepath.Join(projectDir, "main.go")
	args := []string{"run"}

	if _, err := os.Stat(mainFile); err == nil {
		args = append(args, "./"+mainFile)
	} else {
		args = append(args, "./"+projectDir)
	}

	fmt.Printf("🚀 Ejecutando servidor desde %s...\n", mainFile)

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Error al ejecutar el servidor:", err)
	}
}

// detectProjectDir busca automáticamente la carpeta creada por "bwork init"
func detectProjectDir() string {
	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("❌ Error leyendo el directorio actual:", err)
		return ""
	}

	for _, e := range entries {
		if e.IsDir() {
			mainPath := filepath.Join(e.Name(), "main.go")
			if _, err := os.Stat(mainPath); err == nil {
				return e.Name()
			}
		}
	}

	// fallback: si hay main.go en la raíz
	if _, err := os.Stat("main.go"); err == nil {
		return "."
	}

	return ""
}
