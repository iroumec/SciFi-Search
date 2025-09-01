package main

import (
	"fmt"
	"net/http"
)

const (
	port = ":8080" // Puerto que se escucha.
)

var db = setupDB()

func main() {

	registerHandlers()

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
