package main

import (
	"encoding/json"
	"net/http"
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

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
