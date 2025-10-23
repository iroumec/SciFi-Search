package handlers

// ------------------------------------------------------------------------------------------------

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	sqlc "tpe/web/app/database"
	"tpe/web/app/utils"
	"tpe/web/app/views"

	"github.com/a-h/templ"
)

var lastUserID int32 = 12

// ------------------------------------------------------------------------------------------------

// Se registran los endpoints relacionados al manejo de usarios.
func registrarHandlersUsuarios() {

	http.HandleFunc("/users", userHandler)
	http.HandleFunc("/api/users", userHandlerAPI)
	http.HandleFunc("/sign-up", signUpHandler)
}

// ------------------------------------------------------------------------------------------------

func userHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		if utils.HasGETRequestParameters(r) {
			showUserAPI(w, r)
		} else {
			listUsers(w, r)
		}
	case http.MethodPost:
		addUser(w, r)
	case http.MethodPut:
		updateUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

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

// ------------------------------------------------------------------------------------------------

type addUserPayload struct {
	Name       string `json:"name"`
	Middlename string `json:"middlename"`
	Surname    string `json:"surname"`
	// Email    string `json:"email"`    // Descomenta si los usas
	// Password string `json:"password"` // Descomenta si los usas
}

// Agrega un usuario a la base de datos.
func addUser(w http.ResponseWriter, r *http.Request) {

	var payload addUserPayload
	var err error
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if hayCampoIncompleto(payload.Name, payload.Surname) {
		// Campos obligatorios incompletos -> 400 Bad Request
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	// Se publica un evento de creación.
	event := map[string]interface{}{
		"type": "user_created",
		"user": payload,
		"time": time.Now(),
	}

	eventData, _ := json.Marshal(event)
	if err := nat.Publish("products.events", eventData); err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	/*json.NewEncoder(w).Encode(map[string]string{
		"status":  "processing",
		"message": "User creation in progress",
	})*/

	params := sqlc.CreateUserParams{
		UserID:  lastUserID,
		Name:    payload.Name,
		Surname: payload.Surname,
		Middlename: sql.NullString{
			String: payload.Middlename,
			Valid:  payload.Middlename != "", // El campo es 'Valid' si no es un string vacío
		},
	}

	// Creación del usuario en la base de datos.
	_, err = queries.CreateUser(r.Context(), params)
	if err != nil {
		// Error al crear el usuario en la BD -> 500 Internal Server Error.
		log.Printf("Error al crear usuario: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	lastUserID++

	// Se establece el código de estado a 201 Created.
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Acá tendría que haber un "usuario registrado con éxito".
	component := views.SuccessfulSignUpPage()
	templ.Handler(component).ServeHTTP(w, r)
}

// addUserAPI crea un usuario y devuelve el nuevo objeto como JSON.
func addUserAPI(w http.ResponseWriter, r *http.Request) {

	// --- 1. Decodificar y Validar (Sin cambios) ---
	var payload addUserPayload
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

	// --- 4. CREAR EN LA BD (¡Cambio importante!) ---
	//    Para devolver el nuevo objeto, tu consulta sqlc
	//    debe usar `RETURNING *` y ':one'.
	newUser, err := queries.CreateUser(r.Context(), params)
	if err != nil {
		log.Printf("Error al crear usuario: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// --- 5. RESPONDER CON JSON (¡Este es el cambio!) ---

	// Establece el Header a JSON
	w.Header().Set("Content-Type", "application/json")
	// Establece el código de estado a 201 Created
	w.WriteHeader(http.StatusCreated)

	// Codifica el objeto 'newUser' (que incluye el ID de la BD)
	// y lo envía como respuesta JSON.
	json.NewEncoder(w).Encode(newUser)

	// La parte de 'templ' se elimina:
	// component := views.SuccessfulSignUpPage()
	// templ.Handler(component).ServeHTTP(w, r)
}

// ------------------------------------------------------------------------------------------------
// Eliminación de un Usuario
// ------------------------------------------------------------------------------------------------

// Elimina un usuario de la base de datos.
func deleteUser(w http.ResponseWriter, r *http.Request) {

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

	component := views.UserDeletedPage()
	templ.Handler(component).ServeHTTP(w, r)
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

// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------
// Mostrar un Usuario
// ------------------------------------------------------------------------------------------------

// Muestra los datos correspondientes a un usuario, dado un ID.
func showUser(w http.ResponseWriter, r *http.Request) {

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

	component := views.ProfilePage(user)
	templ.Handler(component).ServeHTTP(w, r)
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

// ------------------------------------------------------------------------------------------------

func updateUser(w http.ResponseWriter, r *http.Request) {

	/*
		id, err := extractID(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/

}

func updateUserAPI(w http.ResponseWriter, r *http.Request) {

	/*
		id, err := extractID(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	*/

}

// ------------------------------------------------------------------------------------------------

func extractID(r *http.Request) (int32, error) {

	// Obtención del valor del parámetro 'id' directamente.
	idString := r.URL.Query().Get("id")
	if idString == "" {
		return 0, fmt.Errorf("parámetro 'id' es requerido")
	}

	// Conversión del ID de string a un número, validando que quepa en 32 bits.
	id64, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parámetro 'id' debe ser un número entero válido")
	}

	// Si todo fue exitoso, se convierte el id a int32.
	return int32(id64), nil
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {

	component := views.SignUpPage("")
	templ.Handler(component).ServeHTTP(w, r)
}

// ------------------------------------------------------------------------------------------------
// Listado de Usuarios
// ------------------------------------------------------------------------------------------------

func getListOfUsers(w http.ResponseWriter, r *http.Request) ([]sqlc.User, error) {

	users, err := queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return users, nil
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

// Lista a todos los usuarios.
func listUsers(w http.ResponseWriter, r *http.Request) {

	users, err := queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := views.UserListPage(users)
	templ.Handler(component).ServeHTTP(w, r)
}
