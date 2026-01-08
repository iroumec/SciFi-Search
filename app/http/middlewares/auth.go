package middlewares

// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/http/cookies"
	"scifi-search/app/languages"
)

// ------------------------------------------------------------------------------------------------

func RequiresEmailVerified(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, err := auth.GetCurrentUserID(w, r)
		if err == nil {
			isEmailVerified, err := auth.IsEmailVerified(*userID)
			if err == nil && *isEmailVerified {
				next(w, r)
			}
		}

		message := languages.GetTranslatorFromRequest(r)("Debe verificar su email antes de acceder a esta funcionalidad.")

		// Si es petición HTMX, se envia un trigger para mostrar popup.
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Trigger", `{
						"showFlash": {
						"message": "`+message+`"
						}
					}`)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Si no es HTMX, simplemente se redirige.
		cookies.AddFlashCookie(w, message)
		http.Redirect(w, r, "/", http.StatusFound)
	}
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

		authorizationLevel := auth.GetCurrentAuthorizationLevel(w, r)

		if authorizationLevel < minimumLevelOfAuthorizationRequired {

			message := languages.GetTranslatorFromRequest(r)("No cuenta con los permisos suficientes.")

			// Si es petición HTMX, se envia un trigger para mostrar popup.
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Trigger", `{
					"showFlash": {
						"message": "`+message+`"
					}
				}`)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			// Si no es HTMX, simplemente se redirige.
			cookies.AddFlashCookie(w, message)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next(w, r)
	}
}

// ------------------------------------------------------------------------------------------------
