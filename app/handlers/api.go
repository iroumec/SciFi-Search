package handlers

// ------------------------------------------------------------------------------------------------

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	sqlc "scifi-search/app/database"
	"strconv"
	"strings"
	"time"
)

// ------------------------------------------------------------------------------------------------

func registerAPIHandlers() {
	http.HandleFunc("/api/users", userHandlerAPI)
	http.HandleFunc("/api/users/", userWithIDHandlerAPI)
}

// ------------------------------------------------------------------------------------------------

func userHandlerAPI(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		listUsersAPI(w, r)
	case http.MethodPost:
		addUserAPI(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func userWithIDHandlerAPI(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	id := int32(idInt)

	switch r.Method {
	case http.MethodGet:
		showUserAPI(w, r, id)
	case http.MethodPut:
		updateUserAPI(w, r, id)
	case http.MethodDelete:
		deleteUserAPI(w, r, id)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------
// Creación de Usuario
// ------------------------------------------------------------------------------------------------

// Crea un usuario y devuelve el nuevo objeto como JSON.
func addUserAPI(w http.ResponseWriter, r *http.Request) {

	newUser := addUserToDatabaseAPI(w, r)
	if newUser == nil {
		// Ocurrió un error que ya se trató antes.
		return
	}

	// Se establece el header a JSON.
	w.Header().Set("Content-Type", "application/json")
	// Se establece el código de estado a 201 Created.
	w.WriteHeader(http.StatusCreated)

	// Se codifica el nuevo usuario y se envía como respuesta JSON.
	json.NewEncoder(w).Encode(newUser)
}

// ------------------------------------------------------------------------------------------------
// Eliminación de Usuario
// ------------------------------------------------------------------------------------------------

func deleteUserAPI(w http.ResponseWriter, r *http.Request, id int32) {

	err := queries.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: El usuario no existe.
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			// Error 500: Hubo un problema con la base de datos u otro error inesperado.
			log.Printf("Error al obtener usuario por ID %d: %v", id, err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ------------------------------------------------------------------------------------------------
// Muestra de Usuario
// ------------------------------------------------------------------------------------------------

func showUserAPI(w http.ResponseWriter, r *http.Request, id int32) {

	user, err := queries.GetUserByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: El usuario no existe.
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			// Error 500: Hubo un problema con la base de datos u otro error inesperado.
			log.Printf("Error al obtener usuario por ID %d: %v", id, err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ------------------------------------------------------------------------------------------------
// Actualización de Usuario
// ------------------------------------------------------------------------------------------------

func updateUserAPI(w http.ResponseWriter, r *http.Request, id int32) {

	var payload sqlc.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Cuerpo JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if hayCampoIncompleto(payload.Name, payload.Surname) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	event := map[string]interface{}{
		"type": "user_created",
		"user": payload,
		"time": time.Now(),
	}
	eventData, _ := json.Marshal(event)
	if err := nat.Publish("products.events", eventData); err != nil {
		http.Error(w, "Error procesando la solicitud", http.StatusInternalServerError)
		return
	}

	params := sqlc.UpdateUserParams{
		UserID:  id,
		Name:    payload.Name,
		Surname: payload.Surname,
	}

	err := queries.UpdateUser(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: El usuario no existe.
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			// Error 500: Hubo un problema con la base de datos u otro error inesperado.
			log.Printf("Error al obtener usuario por ID %d: %v", id, err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}
}

// ------------------------------------------------------------------------------------------------
// Listado de Usuarios
// ------------------------------------------------------------------------------------------------

func listUsersAPI(w http.ResponseWriter, r *http.Request) {
	users, err := queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devuelve JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ------------------------------------------------------------------------------------------------

func addUserToDatabaseAPI(w http.ResponseWriter, r *http.Request) *sqlc.User {

	// Se decodifica y valida el payload.
	var payload sqlc.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Cuerpo JSON inválido: "+err.Error(), http.StatusBadRequest)
		return nil
	}

	if hayCampoIncompleto(payload.Name, payload.Surname) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return nil
	}

	// Se publica el evento.
	event := map[string]interface{}{
		"type": "user_created",
		"user": payload,
		"time": time.Now(),
	}
	eventData, _ := json.Marshal(event)
	if err := nat.Publish("products.events", eventData); err != nil {
		http.Error(w, "Error procesando la solicitud", http.StatusInternalServerError)
		return nil
	}

	// Se preparan los parámetros para la BD.
	params := sqlc.CreateUserParams{
		Name:    payload.Name,
		Surname: payload.Surname,
	}

	// Creación del usuario en la base de datos.
	newUser, err := queries.CreateUser(r.Context(), params)
	if err != nil {
		log.Printf("Error al crear usuario: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return nil
	}

	return &newUser
}

// ------------------------------------------------------------------------------------------------
