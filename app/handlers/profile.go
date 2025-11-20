package handlers

// ------------------------------------------------------------------------------------------------
// TODO: a desarrollar en etapas posteriores.
// ------------------------------------------------------------------------------------------------

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	sqlc "scifi-search/app/database"

	_ "github.com/lib/pq"
)

// ------------------------------------------------------------------------------------------------

// Mapea un token de inicio de sesión a un userID.
var sessions = make(map[string]int32)

// ------------------------------------------------------------------------------------------------
// Handle Profile Access
// ------------------------------------------------------------------------------------------------

func handleProfileAccess(user sqlc.User, w http.ResponseWriter, r *http.Request) {

	token, err := generateSessionToken()
	if err != nil {
		http.Error(w, "error interno", http.StatusInternalServerError)
		return
	}

	//sessions[token] = user.ID

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true si se usa HTTPS.
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
