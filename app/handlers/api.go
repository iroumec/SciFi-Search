package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
	sqlc "tpe/web/app/database"
	"tpe/web/app/utils"
)

func registerAPIHandlers() {
	http.HandleFunc("/api/users", userHandlerAPI)
}

func userHandlerAPI(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		if utils.HasGETRequestParameters(r) {
			showUserAPI(w, r)
		} else {
			listUsersAPI(w, r)
		}
	case http.MethodPost:
		addUserAPI(w, r)
	case http.MethodPut:
		updateUserAPI(w, r)
	case http.MethodDelete:
		deleteUserAPI(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// addUserAPI crea un usuario y devuelve el nuevo objeto como JSON.
func addUserAPI(w http.ResponseWriter, r *http.Request) {

	// Se decodifica y valida el payload.
	var payload userPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Cuerpo JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if hayCampoIncompleto(payload.Name, payload.Surname) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
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
		return
	}

	// Se preparan los parámetros para la BD.
	params := sqlc.CreateUserParams{
		UserID:  lastUserID,
		Name:    payload.Name,
		Surname: payload.Surname,
		Middlename: sql.NullString{
			String: payload.Middlename,
			Valid:  payload.Middlename != "",
		},
	}

	lastUserID++

	// Creación del usuario en la base de datos.
	newUser, err := queries.CreateUser(r.Context(), params)
	if err != nil {
		log.Printf("Error al crear usuario: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Se establece el header a JSON.
	w.Header().Set("Content-Type", "application/json")
	// Se establece el código de estado a 201 Created.
	w.WriteHeader(http.StatusCreated)

	// Se codifica el nuevo usuario y se envía como respuesta JSON.
	json.NewEncoder(w).Encode(newUser)
}

func deleteUserAPI(w http.ResponseWriter, r *http.Request) {

	id, err := extractID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = queries.DeleteUser(r.Context(), id)
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

func showUserAPI(w http.ResponseWriter, r *http.Request) {

	id, err := extractID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

func updateUserAPI(w http.ResponseWriter, r *http.Request) {

	id, err := extractID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// --- 1. Decodificar y Validar (Sin cambios) ---
	var payload userPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Cuerpo JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	if hayCampoIncompleto(payload.Name, payload.Surname) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	// --- 2. Publicar Evento (Sin cambios) ---
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

	// --- 3. Preparar Parámetros de BD (Sin cambios) ---
	//    (He eliminado la lógica de 'lastUserID' porque la BD
	//    debería generar el ID automáticamente, por ejemplo, con SERIAL)
	params := sqlc.UpdateUserParams{
		UserID:  id,
		Name:    payload.Name,
		Surname: payload.Surname,
		Middlename: sql.NullString{
			String: payload.Middlename,
			Valid:  payload.Middlename != "",
		},
	}

	err = queries.UpdateUser(r.Context(), params)
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

// listUsersAPI devuelve la lista de usuarios como un array JSON.
func listUsersAPI(w http.ResponseWriter, r *http.Request) {
	users, err := queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devuelve JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
