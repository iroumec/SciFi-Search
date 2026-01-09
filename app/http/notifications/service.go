package notifications

// ------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/http/notifications/cookies"
	"scifi-search/app/http/notifications/triggers"
)

// ------------------------------------------------------------------

// Retorna el tipo de Notificación que se agregó.
func ShowFlash(w http.ResponseWriter, r *http.Request, message string) Notifications {

	// Si es petición HTMX, se envia un trigger para mostrar popup.
	if triggerAdded := triggers.AddHXTrigger(w, r, message); triggerAdded {
		return HXTriggerNotification
	}

	// Si no es HTMX, simplemente se redirige.
	cookies.AddFlashCookie(w, message)
	return CookieNotification
}

// ------------------------------------------------------------------
