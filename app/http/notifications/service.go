package notifications

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/http/notifications/cookies"
	"scifi-search/app/http/notifications/triggers"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Returns the type of notification added.
func ShowFlash(w http.ResponseWriter, r *http.Request, message string) Notifications {

	// HTMX -> Popup activated via trigger.
	if triggerAdded := triggers.AddHXTrigger(w, r, message); triggerAdded {
		return HXTriggerNotification
	}

	// No HTMX -> Cookie.
	cookies.AddFlashCookie(w, message)
	return CookieNotification
}

// ------------------------------------------------------------------
