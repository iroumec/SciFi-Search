package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import "net/http"

// ------------------------------------------------------------------------------------------------
// Constantes
// ------------------------------------------------------------------------------------------------

const (
	browserSimpleMessage = "OK"
)

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

func RegisterHealth() {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, MethodNotAllowedError.Error(), http.StatusMethodNotAllowed)
			return
		}

		// Status Code: 200 OK.
		// This is request by `curl -f` when the server is initiated.
		w.WriteHeader(http.StatusOK)

		// Simple body in case the page is opened in a browser.
		w.Write([]byte(browserSimpleMessage))
	})
}

// ------------------------------------------------------------------------------------------------
