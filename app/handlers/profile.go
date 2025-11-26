package handlers

// ------------------------------------------------------------------------------------------------
// TODO: a desarrollar en etapas posteriores.
// ------------------------------------------------------------------------------------------------

import (
	"net/http"

	"scifi-search/app/views"

	_ "github.com/lib/pq"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
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

	// Si el usuario está autenticado (tiene cookies de sesión)...
	if isUserAuthenticated(r) {

		sessionContainer, _ := session.GetSession(r, nil, &sessmodels.VerifySessionOptions{
			SessionRequired: boolPtr(false),
		})

		supertokensUserID := sessionContainer.GetUserID()

		user, err := queries.GetUserByAuthID(r.Context(), supertokensUserID)
		if err != nil {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}

		component := views.LoggedProfilePage(user)
		component.Render(r.Context(), w)

	} else {

		component := views.UnloggedProfilePage()
		component.Render(r.Context(), w)
	}
}
