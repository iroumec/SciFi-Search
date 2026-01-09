package cookies

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"encoding/base64"
	"net/http"
)

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

// Añade una cookie flash (de muy corta duración).
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
