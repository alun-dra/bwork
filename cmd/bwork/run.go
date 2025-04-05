package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runServer() {
	fmt.Println("🚀 Ejecutando servidor desde app/main.go...")
	cmd := exec.Command("go", "run", "app/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error al ejecutar el servidor:", err)
	}
}
