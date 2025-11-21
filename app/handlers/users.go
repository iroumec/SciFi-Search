package handlers

// ------------------------------------------------------------------------------------------------

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"scifi-search/app/utils"
	"scifi-search/app/views"
	"strconv"
	"strings"
	"time"

	sqlc "scifi-search/app/database"

	"github.com/a-h/templ"
)

// ------------------------------------------------------------------------------------------------

// Se registran los endpoints relacionados al manejo de usarios.
func registrarHandlersUsuarios() {
	http.HandleFunc("/users", userHandler)
	http.HandleFunc("/users/", userWithIDHandler)
	registerAPIHandlers()
}

// ------------------------------------------------------------------------------------------------

func userWithIDHandler(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	id := int32(idInt)

	switch r.Method {
	case http.MethodDelete:
		deleteUser(w, r, id)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func userHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		if utils.HasGETRequestParameters(r) {
			showUser(w, r)
		} else {
			listUsers(w, r)
		}
	case http.MethodPost:
		addUser(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------
// Agregado de un Usuario.
// ------------------------------------------------------------------------------------------------

// Agrega un usuario a la base de datos.
func addUser(w http.ResponseWriter, r *http.Request) {

	newUser := addUserToDatabase(w, r)
	if newUser == nil {
		return
	}

	component := views.UserIndividual(*newUser)
	templ.Handler(component).ServeHTTP(w, r)
}

// ------------------------------------------------------------------------------------------------
// Eliminación de un Usuario
// ------------------------------------------------------------------------------------------------

// Elimina un usuario de la base de datos.
func deleteUser(w http.ResponseWriter, r *http.Request, id int32) {

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

	log.Printf("Se eliminó al usuario de ID %d.", id)

	// Por defecto, la respuesta es 200 OK.
}

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

	component := views.LoggedProfilePage(user)
	templ.Handler(component).ServeHTTP(w, r)
}

// ------------------------------------------------------------------------------------------------
// Extracción de ID de la URL
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

// ------------------------------------------------------------------------------------------------

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

// ------------------------------------------------------------------------------------------------
// Agregado de Usuario a la Base de Datos
// ------------------------------------------------------------------------------------------------

func addUserToDatabase(w http.ResponseWriter, r *http.Request) *sqlc.User {

	// Parseo del formulario enviado por POST.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario: "+err.Error(), http.StatusBadRequest)
		return nil
	}

	// Construcción del payload.
	var payload sqlc.User
	payload.Name = r.Form.Get("name")
	payload.Surname = r.Form.Get("surname")

	// Validación.
	if utils.HayCampoIncompleto(payload.Name, payload.Surname) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return nil
	}

	// Publicación de un evento.
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

	// Parámetros para la DB.
	params := sqlc.CreateUserParams{
		Name:    payload.Name,
		Surname: payload.Surname,
	}

	// Creación en la base.
	newUser, err := queries.CreateUser(r.Context(), params)
	if err != nil {
		log.Printf("Error al crear usuario: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return nil
	}

	return &newUser
}
