package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"slices"

	"golang.org/x/crypto/bcrypt"

	sqlc "uki/app/database/sqlc"

	_ "github.com/lib/pq"
)

// registerHandlers registra todos los endpoints
func registerUserHandlers() {

	fmt.Println("Registrando handlers de usuarios...")

	// Handler que maneja el registro de usuarios.
	http.HandleFunc("/signin", signInHandler)

	// Handler que maneja el login de usuarios.
	http.HandleFunc("/login", logInHandler)

	fmt.Println("Handlers de usuarios registrados...")
}

// ------------------------------------------------------------------------------------------------
// SignIn Handler
// ------------------------------------------------------------------------------------------------

func signInHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		signInHandleGET(w)
	case http.MethodPost:
		signInHandlePOST(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func signInHandleGET(w http.ResponseWriter) {

	renderizeTemplate(w, "template/signin.html", nil)
}

// ------------------------------------------------------------------------------------------------

func signInHandlePOST(w http.ResponseWriter, r *http.Request) {

	// Se parsean los datos del formulario enviados vía POST.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if hayCampoIncompleto(username, name, email, password) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	// Se encripta la contraseña para no manejar credenciales en bruto.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	createdUser, err := queries.CreateUser(r.Context(), sqlc.CreateUserParams{
		Username: username,
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		log.Printf("error creating user: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	// Quizás pueda usarse después...
	_ = createdUser

	// Se redirige al usuario al menú princiapl.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------
// LogIn Handler
// ------------------------------------------------------------------------------------------------

func logInHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		logInHandleGET(w, "")
	case http.MethodPost:
		logInHandlePOST(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func logInHandleGET(w http.ResponseWriter, errorMessage string) {

	data := map[string]any{
		"ErrorMessage": errorMessage,
	}

	renderizeTemplate(w, "template/login.html", data)
}

// ------------------------------------------------------------------------------------------------

func logInHandlePOST(w http.ResponseWriter, r *http.Request) {

	// Se parsean los datos del formulario enviados vía POST.
	if err := r.ParseForm(); err != nil {
		// Se renderiza la página con el error correspondiente.
		logInHandleGET(w, "Error al procesar el formulario.")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if hayCampoIncompleto(username, password) {
		logInHandleGET(w, "Faltan campos obligatorios.")
		return
	}

	user, err := queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			logInHandleGET(w, "El usuario proporcionado no existe.")
			return
		}
		log.Printf("error getting user: %v", err)
		logInHandleGET(w, "Error interno del servidor.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logInHandleGET(w, "Contraseña incorrecta.")
		return
	}

	handleProfileAccess(user, w, r)
}

// ------------------------------------------------------------------------------------------------
// Verificación de campos
// ------------------------------------------------------------------------------------------------

func hayCampoIncompleto(campos ...string) bool {

	return slices.Contains(campos, "")
}

// ------------------------------------------------------------------------------------------------
// Parseo de datos
// ------------------------------------------------------------------------------------------------

func parseData(campos ...string) bool {

	return slices.Contains(campos, "")
}
