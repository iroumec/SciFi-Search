package handlers

// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/languages"
	"scifi-search/app/utils/structures"
	"scifi-search/app/views"
)

// ------------------------------------------------------------------------------------------------

// Registra los endpoint asociados al historial.
func registerHistoryHandlers() {

	http.HandleFunc("/history", historyHandler)

}

// ------------------------------------------------------------------------------------------------

func historyHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		showHistory(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}

}

// ------------------------------------------------------------------------------------------------

func showHistory(w http.ResponseWriter, r *http.Request) {
	user := getCurrentUser(w, r)

	// Si hay un usuario autenticado (con cookies de sesión)...
	if user != nil {

		rows, err := queries.ListHistoricSearchesFromUser(r.Context(), user.UserID)
		if err != nil {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		searches := make([]structures.HistoricSearchView, len(rows))
		for i, r := range rows {
			searches[i] = structures.HistoricSearchView{
				Query: r.SearchString,
				Date:  r.SearchDatetime,
			}
		}

		component := views.HistoryPage(searches, languages.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)

	} else {
		component := views.UnloggedPage(languages.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)
	}
}
