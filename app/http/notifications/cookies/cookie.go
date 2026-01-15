package cookies

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"encoding/base64"
	"net/http"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Adds a flash cookie (short duration).
func AddFlashCookie(w http.ResponseWriter, message string) {

	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    base64.StdEncoding.EncodeToString([]byte(message)),
		Path:     "/", // Allows the cookie to be available across all the application.
		MaxAge:   10,
		HttpOnly: false,
	})
}

// ------------------------------------------------------------------------------------------------
