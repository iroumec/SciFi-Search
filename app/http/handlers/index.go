package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/languages"
	"scifi-search/app/views"
)

// ------------------------------------------------------------------------------------------------
// Servicios
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
