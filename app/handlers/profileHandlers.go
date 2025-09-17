package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	sqlc "uki/app/database/sqlc"

	_ "github.com/lib/pq"
)

// Mapea un token de inicio de sesión a un userID.
var sessions = make(map[string]int32)

// ------------------------------------------------------------------------------------------------
// Register Profile Handlers
// ------------------------------------------------------------------------------------------------

// Registra todos los endpoint relacionados al perfil de usuario.
func registerProfileHandlers() {

	fmt.Println("Registrando handlers de perfil...")

	// Handler que maneja el acceso al perfil.
	http.HandleFunc("/profile", profileHandler)

	fmt.Println("Handlers de perfil registrados...")
}

// ------------------------------------------------------------------------------------------------
// Profile Handler
// ------------------------------------------------------------------------------------------------

func profileHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("\nManejando renderizado del perfil...")

	tmpl := template.Must(template.ParseFiles("template/profile.html"))

	c, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "no autenticado", http.StatusUnauthorized)
		return
	}

	userID, ok := sessions[c.Value]
	if !ok {
		http.Error(w, "sesión inválida", http.StatusUnauthorized)
		return
	}

	user, err := queries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "error cargando usuario", http.StatusInternalServerError)
		return
	}

	fmt.Println("Renderizando plantilla de acuerdo a los datos del usuario...")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, user)
	if err != nil {
		http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
	}

	fmt.Println("Plantilla renderizada.")
}

// ------------------------------------------------------------------------------------------------
// Handle Profile Access
// ------------------------------------------------------------------------------------------------

func handleProfileAccess(user sqlc.User, w http.ResponseWriter, r *http.Request) {

	token, err := generateSessionToken()
	if err != nil {
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	sessions[token] = user.ID

	fmt.Println("Se agregó un nuevo token de sesión.")

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true si se usa HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	// Se redirige al usuario a la página del perfil.
	// F12 -> Network. Debería verse un 303 si esto funciona bien.
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------
// Generate Session Token
// ------------------------------------------------------------------------------------------------

// Genera un token aleatoria para el inicio de sesión.
func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
