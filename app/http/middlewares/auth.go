package middlewares

import (
	"log"
	"net/http"
	"scifi-search/app/http/cookies"
	"scifi-search/app/utils/converters"
	"slices"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
)

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Obtención manual de la sesión. Esto la verifica y la pone en el contexto.
		sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
			SessionRequired: converters.ToBoolPointer(true),
		})

		if err != nil {
			log.Printf("Error getting session: %v", err)
			http.Error(w, "Unauthorized - No valid session", http.StatusUnauthorized)
			return
		}

		if sessionContainer == nil {
			log.Printf("Session container is nil")
			http.Error(w, "Unauthorized - No session", http.StatusUnauthorized)
			return
		}

		userID := sessionContainer.GetUserID()

		roles, err := userroles.GetRolesForUser("public", userID, nil)
		if err != nil || !slices.Contains(roles.OK.Roles, "admin") {
			cookies.AddFlashCookie(w, "No cuenta con los permisos suficientes.")
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
