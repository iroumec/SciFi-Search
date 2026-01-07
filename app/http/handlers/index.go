package handlers

import (
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/languages"
	"scifi-search/app/views"
)

func RegisterIndex() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Se crea una instancia de la componente de página.
		component := views.IndexPage(auth.GetCurrentAuthorizationLevel(w, r, queries), languages.GetTranslatorFromRequest(r))

		// Se renderiza la componente.
		component.Render(r.Context(), w)
	})
}
