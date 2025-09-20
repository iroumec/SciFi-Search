package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"golang.org/x/crypto/bcrypt"

	sqlc "uki/app/database"

	_ "github.com/lib/pq"
)

// registerHandlers registra todos los endpoints
func registerUserHandlers() {

	// Handler que maneja el registro de usuarios.
	http.HandleFunc("/signin", signInHandler)

	// Handler que maneja el login de usuarios.
	http.HandleFunc("/login", logInHandler)
}

// ------------------------------------------------------------------------------------------------
// SignIn Handler
// ------------------------------------------------------------------------------------------------

func signInHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		signInHandleGET(w, "")
	case http.MethodPost:
		signInHandlePOST(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func signInHandleGET(w http.ResponseWriter, errorMessage string) {

	data := map[string]any{
		"ErrorMessage": errorMessage,
	}

	renderizeTemplate(w, "template/usuarios/registro/registrarse.html", data, nil)
}

// ------------------------------------------------------------------------------------------------

func signInHandlePOST(w http.ResponseWriter, r *http.Request) {

	// Se parsean los datos del formulario enviados vía POST.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	dni := r.FormValue("dni")
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if hayCampoIncompleto(dni, name, email, password) {
		signInHandleGET(w, "Faltan campos obligatorios.")
		return
	}

	// Se encripta la contraseña para no manejar credenciales en bruto.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	createdUser, err := queries.CrearUsuario(r.Context(), sqlc.CrearUsuarioParams{
		Dni:        dni,
		Nombre:     name,
		Email:      email,
		Contraseña: string(hashedPassword),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			http.Error(w, "El usuario o email ya existen.", http.StatusConflict)
			return
		}

		log.Printf("error creando usuario: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	// Quizás pueda usarse después...
	_ = createdUser

	renderizeTemplate(w, "template/usuarios/registro/registro-exitoso.html", nil, nil)
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

	renderizeTemplate(w, "template/usuarios/iniciar-sesion.html", data, nil)
}

// ------------------------------------------------------------------------------------------------

func logInHandlePOST(w http.ResponseWriter, r *http.Request) {

	// Se parsean los datos del formulario enviados vía POST.
	if err := r.ParseForm(); err != nil {
		// Se renderiza la página con el error correspondiente.
		logInHandleGET(w, "Error al procesar el formulario.")
		return
	}

	dni := r.FormValue("dni")
	password := r.FormValue("password")

	if hayCampoIncompleto(dni, password) {
		logInHandleGET(w, "Faltan campos obligatorios.")
		return
	}

	// Hacer funcionar esto...
	user, err := queries.ObtenerUsuarioPorDNI(r.Context(), dni)
	if err != nil {
		if err == sql.ErrNoRows {
			logInHandleGET(w, "El usuario proporcionado no existe.")
			return
		}
		log.Printf("error getting user: %v", err)
		logInHandleGET(w, "Error interno del servidor.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Contraseña), []byte(password)); err != nil {
		logInHandleGET(w, "Contraseña incorrecta.")
		return
	}

	handleProfileAccess(user, w, r)
}
