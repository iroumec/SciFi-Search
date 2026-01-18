package cookies

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"encoding/base64"
	"net/http"
	"net/url"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Adds a flash cookie (short duration).
func AddFlashCookie(w http.ResponseWriter, message string) {
	encodedMessage := base64.StdEncoding.EncodeToString([]byte(message))
	escapedMessage := url.QueryEscape(encodedMessage)

	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    escapedMessage,
		Path:     "/", // Allows the cookie to be available across all the application.
		MaxAge:   10,
		HttpOnly: false,
	})
}

// ------------------------------------------------------------------------------------------------
