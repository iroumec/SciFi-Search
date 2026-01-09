package middlewares

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/http/notifications"
	"scifi-search/app/languages"
)

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

func RequiresEmailVerified(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, err := auth.GetCurrentUserID(w, r)
		if err == nil {
			isEmailVerified, err := auth.IsEmailVerified(*userID)
			if err == nil && *isEmailVerified {
				next(w, r)
				return
			}
		}

		message := languages.GetTranslatorFromRequest(r)("Debe verificar su email antes de acceder a esta funcionalidad.")

		addFlash(w, r, message)
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

			addFlash(w, r, message)
			return
		}
		next(w, r)
	}
}

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

func addFlash(w http.ResponseWriter, r *http.Request, message string) {

	notificationType := notifications.ShowFlash(w, r, message)

	switch notificationType {
	case notifications.HXTriggerNotification:
		w.WriteHeader(http.StatusForbidden)
	case notifications.CookieNotification:
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}
}

// ------------------------------------------------------------------------------------------------
