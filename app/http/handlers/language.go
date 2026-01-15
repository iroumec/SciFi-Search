package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"errors"
	"net/http"
)

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

var (
	InvalidLanguageError = errors.New("error.invalid-language")
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Endpoints.
func RegisterLanguageHandlers() {

	http.HandleFunc("/language", languageHandler)
}

// ------------------------------------------------------------------------------------------------
// Functions
// ------------------------------------------------------------------------------------------------

func languageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowedError.Error(), http.StatusBadRequest)
		return
	}

	err := setLanguage(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	goBackToPreviousPage(w, r)
}

// ------------------------------------------------------------------------------------------------

func setLanguage(w http.ResponseWriter, r *http.Request) error {

	language := r.URL.Query().Get("language")

	// Language validation.
	if language != "es" && language != "en" {
		return InvalidLanguageError
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "language",
		Value: language,
		Path:  "/",
	})

	return nil
}

// ------------------------------------------------------------------------------------------------

func goBackToPreviousPage(w http.ResponseWriter, r *http.Request) {

	referer := r.Header.Get("Referer")

	// Fallback.
	if referer == "" {
		referer = "/"
	}

	w.Header().Set("HX-Redirect", referer)
	w.WriteHeader(http.StatusNoContent)
}

// ------------------------------------------------------------------------------------------------
