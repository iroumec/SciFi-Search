package main

import (
	"fmt"
	"log"
	"net/http"

	sqlc "uki/code/database/sqlc"

	_ "github.com/lib/pq"
)

const (
	fileDir = "../static" // Directorio relativo con los archivos estáticos. Relativo adonde se ejecuta go run.
)

// registerHandlers registra todos los endpoints
func registerHandlers() {

	// Se crea un manejador (handler) de servidor de archivos.
	fileServer := http.FileServer(http.Dir(fileDir))

	// Se envuelve en un gzip middleware.
	http.Handle("/", gzipMiddleware(fileServer))

	// Se define un handler que maneje la creación de usuarios.
	http.HandleFunc("/users", usersHandler)
}

// usersHandler maneja creación de usuarios
func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")

	if name == "" || email == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	createdUser, err := queries.CreateUser(ctx,
		sqlc.CreateUserParams{
			Name:  name,
			Email: email,
		})

	if err != nil {
		log.Fatal("error creating user:", err)
	}

	fmt.Printf("User created: %+v\n", createdUser)
}
