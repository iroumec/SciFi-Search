package handlers

// ------------------------------------------------------------------------------------------------
// TODO: a desarrollar en etapas posteriores.
// ------------------------------------------------------------------------------------------------

import (
	"net/http"

	"scifi-search/app/views"

	_ "github.com/lib/pq"
)

// ------------------------------------------------------------------------------------------------

func registerProfileHandlers() {

	http.HandleFunc("/profile", handleProfile)
}

// ------------------------------------------------------------------------------------------------

func handleProfile(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		showProfile(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func showProfile(w http.ResponseWriter, r *http.Request) {

	user := getCurrentUser(w, r)

	// Si hay un usuario autenticado (con cookies de sesión)...
	if user != nil {

		searches, err := queries.ListHistoricSearchesFromUser(r.Context(), user.UserID)
		if err != nil {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}

		component := views.LoggedProfilePage(*user, searches)
		component.Render(r.Context(), w)

	} else {

		component := views.UnloggedProfilePage()
		component.Render(r.Context(), w)
	}
}
