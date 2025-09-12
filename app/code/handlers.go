package code

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	sqlc "uki/app/database/sqlc"

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

	defer r.Body.Close()

	var req struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "request body inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Name == "" || req.Email == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	createdUser, err := queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Username: req.Username,
		Name:     req.Name,
		Email:    req.Email,
	})
	if err != nil {
		log.Printf("error creating user: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"ID": createdUser.ID})

	fmt.Printf("User created: %+v\n", createdUser)
}
