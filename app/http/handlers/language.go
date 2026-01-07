package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"net/http"

	"scifi-search/app/languages"
)

// ------------------------------------------------------------------------------------------------
// Registro de endpoints
// ------------------------------------------------------------------------------------------------

// Carga los mensajes y registra los handlers correspondientes al lenguaje.
func RegisterLanguageHandlers() {

	languages.LoadAllMessages()

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
		http.Error(w, "invalid language", http.StatusBadRequest)
		return
	}

	// Seteo de cookie.
	http.SetCookie(w, &http.Cookie{
		Name:  "language",
		Value: language,
		Path:  "/",
	})

	// Volver a la página previa.
	referer := r.Header.Get("Referer")
	if referer != "" {
		http.Redirect(w, r, referer, http.StatusSeeOther)
		return
	}

	// Si no hay referer, fallback a página principal.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------
