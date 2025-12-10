package handlers

import (
	"net/http"
	"scifi-search/app/utils"
	"scifi-search/app/views"
)

func registerSettingsHandlers() {

	http.HandleFunc("/settings", handleSettings)

}

func handleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		showSettings(w, r)
	default: 
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func showSettings(w http.ResponseWriter, r *http.Request) {

	if getCurrentUser(w,r) != nil {

		component := views.SettingsPage(true, utils.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)

	} else {

		component := views.SettingsPageError()
		component.Render(r.Context(), w)

	}

}