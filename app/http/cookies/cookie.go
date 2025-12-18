package cookies

// ------------------------------------------------------------------------------------------------

import (
	"encoding/base64"
	"net/http"
)

// ------------------------------------------------------------------------------------------------

func AddFlashCookie(w http.ResponseWriter, message string) {

	http.SetCookie(w, &http.Cookie{
		Name:     "flash",
		Value:    base64.StdEncoding.EncodeToString([]byte(message)),
		Path:     "/", // Para que la cookie esté disponible en toda la app.
		MaxAge:   10,
		HttpOnly: false,
	})
}

// ------------------------------------------------------------------------------------------------
