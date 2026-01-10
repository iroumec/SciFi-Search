package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"errors"
	"net/http"
)

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

var (
	InvalidLanguageError = errors.New("Invalid language")
)

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

// Registro de endpoints.
func RegisterLanguageHandlers() {

	http.HandleFunc("/language", setLanguageHandler)
}

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

// Establece el lenguaje de la aplicación.
func setLanguageHandler(w http.ResponseWriter, r *http.Request) {

	// Obtención del lenguaje de la URL.
	language := r.URL.Query().Get("language")

	// Validación del lenguaje.
	if language != "es" && language != "en" {
		http.Error(w, InvalidLanguageError.Error(), http.StatusBadRequest)
		return
	}

	// Seteo de cookie.
	http.SetCookie(w, &http.Cookie{
		Name:  "language",
		Value: language,
		Path:  "/",
	})

	// Se vuelve a la página previa.
	referer := r.Header.Get("Referer")
	if referer != "" {
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}

	// Si no hay referer, fallback a página principal.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------
