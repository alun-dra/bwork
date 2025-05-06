package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runServer() {
	moduleName := getModuleName()
	fmt.Printf("🚀 Ejecutando servidor desde %s/main.go...\n", moduleName)

	cmd := exec.Command("go", "run", "./"+moduleName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Error al ejecutar el servidor:", err)
	}
}
