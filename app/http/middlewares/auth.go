package middlewares

// ------------------------------------------------------------------------------------------------

import (
	"log"
	"net/http"
	"scifi-search/app/http/cookies"
	"scifi-search/app/infra/auth"
	"scifi-search/app/languages"
	"scifi-search/app/utils/converters"

	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

// ------------------------------------------------------------------------------------------------

func RequiresEmailVerified(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isEmailVerified(w, r) {

			message := languages.GetTranslatorFromRequest(r)("Debe verificar su email antes de acceder a esta funcionalidad.")
			cookies.AddFlashCookie(w, message)

			// Si es petición HTMX, se envia un trigger para mostrar popup.
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Trigger", "showFlash")
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Si no es HTMX, simplemente se redirige.
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next(w, r)
	}
}

// Función auxiliar para verificar si el email está verificado.
func isEmailVerified(w http.ResponseWriter, r *http.Request) bool {
	sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
		SessionRequired: converters.ToBoolPointer(true),
	})
	if err != nil {
		return false
	}

	if sessionContainer == nil {
		return false
	}

	userID := sessionContainer.GetUserID()
	isVerified, err := emailverification.IsEmailVerified(userID, nil, nil)

	return isVerified
}

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

			message := languages.GetTranslatorFromRequest(r)("No cuenta con los permisos suficientes.")
			cookies.AddFlashCookie(w, message)

			// Si es petición HTMX, se envia un trigger para mostrar popup.
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Trigger", "showFlash")
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Si no es HTMX, simplemente se redirige.
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next(w, r)
	}
}

// ------------------------------------------------------------------------------------------------
