package middlewares

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/http/notifications"
	"scifi-search/app/languages"
)

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

// Errors
var (
	EmailNotVerifiedSuggestion     = "suggestion.email-not-verified"
	NotEnoughPermissionsSuggestion = "suggestion.not-enough-permissions"
)

// ------------------------------------------------------------------------------------------------
// Services
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

		message := languages.GetTranslatorFromRequest(r)(EmailNotVerifiedSuggestion)

		addFlash(w, r, message)
	}
}

// ------------------------------------------------------------------------------------------------

// Middleware HTTP that protects endpoints according to the authorization level of the user.
//
// Verifys that a valid session exists and the hierarchical user authorization level.
func RequiresAuthorization(
	next http.HandlerFunc, minimumLevelOfAuthorizationRequired int,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authorizationLevel := auth.GetCurrentAuthorizationLevel(w, r)

		if authorizationLevel < minimumLevelOfAuthorizationRequired {

			message := languages.GetTranslatorFromRequest(r)(NotEnoughPermissionsSuggestion)

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
