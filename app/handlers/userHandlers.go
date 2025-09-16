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
	http.HandleFunc("/signIn", signInHandler)

	// Handler que maneja el login de usuarios.
	http.HandleFunc("/logIn", logInHandler)

	fmt.Println("Handlers de usuarios registrados...")
}

// ------------------------------------------------------------------------------------------------
// SignIn Handler
// ------------------------------------------------------------------------------------------------

func signInHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Print("Manejando registro de usuario...")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

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

	fmt.Println("Manejando log in de usuario...")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Se parsean los datos del formulario enviados vía POST.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if hayCampoIncompleto(username, password) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	user, err := queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "El usuario proporcionado no existe.", http.StatusNotFound)
			return
		}
		log.Printf("error getting user: %v", err)
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Contraseña incorrecta.", http.StatusUnauthorized)
		return
	}

	fmt.Println("Login extisoso, generando sesión...")

	handleProfileAccess(user, w, r)
}

// ------------------------------------------------------------------------------------------------
// Verificación de campos
// ------------------------------------------------------------------------------------------------

func hayCampoIncompleto(campos ...string) bool {

	return slices.Contains(campos, "")
}
