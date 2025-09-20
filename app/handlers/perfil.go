package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	sqlc "uki/app/database"

	_ "github.com/lib/pq"
)

// Mapea un token de inicio de sesión a un userID.
var sessions = make(map[string]int32)

// ------------------------------------------------------------------------------------------------
// Profile Handler
// ------------------------------------------------------------------------------------------------

func manejarPerfil(w http.ResponseWriter, r *http.Request) {

	fmt.Println("\nManejando renderizado del perfil...")

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

	user, err := queries.ObtenerUsuarioPorID(r.Context(), userID)
	if err != nil {
		http.Error(w, "error cargando usuario", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Name":     user.Dni,
		"Username": user.Nombre,
		"Email":    user.Email,
	}

	renderizeTemplate(w, "template/profile.html", data, nil)

	fmt.Println("Plantilla renderizada.")
}

// ------------------------------------------------------------------------------------------------
// Handle Profile Access
// ------------------------------------------------------------------------------------------------

func handleProfileAccess(user sqlc.Usuario, w http.ResponseWriter, r *http.Request) {

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
