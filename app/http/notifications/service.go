package notifications

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/http/notifications/cookies"
	"scifi-search/app/http/notifications/triggers"
)

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

// Retorna el tipo de notificación que se agregó.
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
