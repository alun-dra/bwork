package main

import (
    "fmt"
    "log"
    "net/http"

    "app"
)

func runServer() {
    mux := http.NewServeMux()
    app.SetupRoutes(mux)

    fmt.Println("Servidor escuchando en http://localhost:8081 ...")
    err := http.ListenAndServe(":8081", mux)
    if err != nil {
        log.Fatal("Error al iniciar el servidor:", err)
    }
}
