package handlers

// ------------------------------------------------------------------------------------------------

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	sqlc "tpe/web/app/database"
	"tpe/web/app/utils"
	"tpe/web/app/views"

	"github.com/a-h/templ"
)

// ------------------------------------------------------------------------------------------------

// Se registran los endpoints relacionados al manejo de usarios.
func registrarHandlersUsuarios() {

	http.HandleFunc("/users", userHandler)
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
	case http.MethodPut:
		updateUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

// Agrega un usuario a la base de datos.
func addUser(w http.ResponseWriter, r *http.Request) {

	// Se parsea el body de la request.
	err := r.ParseForm()
	if err != nil {
		// Cuerpo mal formado
	}

	name := r.FormValue("name")
	middlename := r.FormValue("middlename")
	surname := r.FormValue("surname")

	if hayCampoIncompleto(name, surname) {
		// Campos obligatorios incompletos -> 400 Bad Request
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}
	_, err = queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Name:       name,
		Middlename: sql.NullString{String: middlename, Valid: middlename != ""},
		Surname:    surname,
	})
	if err != nil {
		// Error al crear el usuario en la BD -> 500 Internal Server Error.
		log.Printf("Error al crear usuario: %v", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Se establece el código de estado a 201 Created.
	w.WriteHeader(http.StatusCreated)

	// Acá tendría que haber un "usuario registrado con éxito".
	component := views.SuccessfulSignUpPage()
	templ.Handler(component).ServeHTTP(w, r)
}

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

// ------------------------------------------------------------------------------------------------

// Lista a todos los usuarios.
func listUsers(w http.ResponseWriter, r *http.Request) {

	users, err := queries.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := views.UserPage(users)
	templ.Handler(component).ServeHTTP(w, r)
}

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
