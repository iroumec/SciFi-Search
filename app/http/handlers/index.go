package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/languages"
	"scifi-search/app/views"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

func RegisterIndex() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		component := views.IndexPage(
			auth.GetCurrentAuthorizationLevel(w, r),
			languages.GetTranslatorFromRequest(r),
		)
		component.Render(r.Context(), w)
	})
}

// ------------------------------------------------------------------------------------------------
