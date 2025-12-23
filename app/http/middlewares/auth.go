package middlewares

// ------------------------------------------------------------------------------------------------

import (
	"log"
	"net/http"
	"scifi-search/app/http/cookies"
	"scifi-search/app/infra/auth"
	"scifi-search/app/utils/converters"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

// ------------------------------------------------------------------------------------------------

// RequiresAuthorization es un middleware HTTP que protege endpoints según el nivel
// de autorización del usuario.
//
// El middleware verifica que exista una sesión válida y que el usuario tenga al menos
// el nivel de autorización especificado. Los niveles son jerárquicos:
//
//   - 2 (admin): Acceso completo al sistema
//   - 1 (loader): Puede crear y gestionar sus propios recursos
//   - 0 (user): Solo lectura
//
// Si la verificación falla, retorna:
//   - 401 Unauthorized: No hay sesión válida
//   - 403 Forbidden: El usuario no tiene el nivel requerido
//
// Ejemplo de uso:
//
//	http.HandleFunc("/admin/users", RequiresAuthorization(usersHandler, 1))
func RequiresAuthorization(next http.HandlerFunc, minimumLevelOfAuthorizationRequired int) http.HandlerFunc {
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
		authenticationLevel := auth.GetAuthenticationLevel(userID)

		if authenticationLevel < minimumLevelOfAuthorizationRequired {
			cookies.AddFlashCookie(w, "No cuenta con los permisos suficientes.")
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

// ------------------------------------------------------------------------------------------------
