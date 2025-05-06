package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "init":
		moduleName := "app"
		if len(os.Args) >= 3 {
			moduleName = os.Args[2]
		}
		runInit(moduleName)
	case "install":
		if len(os.Args) < 3 {
			fmt.Println("Uso: bwork install <lib>")
		} else {
			runInstall(os.Args[2])
		}
	case "generate":
		if len(os.Args) < 4 {
			fmt.Println("Uso: bwork generate <tipo> <nombre>")
		} else {
			moduleName := getModuleName() // <- debes implementar esta función
			runGenerate(os.Args[2], os.Args[3], moduleName)
		}

	case "run":
		runServer()
	case "help":
		printHelp()
	default:
		fmt.Println("Comando no reconocido:", os.Args[1])
		printHelp()
	}
}

func getModuleName() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Println("❌ No se pudo leer go.mod:", err)
		return "app"
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}

	// Valor por defecto si no se encuentra la línea "module"
	return "app"
}

func printHelp() {
	fmt.Println("Comandos disponibles:")
	fmt.Println("  init                         Inicializa un nuevo proyecto")
	fmt.Println("  install <lib>                Instala una librería local")
	fmt.Println("  generate controller <name>  Crea un controlador")
	fmt.Println("  generate model <name>       Crea un modelo")
	fmt.Println("  generate view <name>        Crea una vista (endpoint)")
	fmt.Println("  generate module <name>        Crea un modulo completo con controller, model, view")
	fmt.Println("  run                          Ejecuta el servidor en puerto 8081")
}
