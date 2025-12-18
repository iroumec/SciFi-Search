package handlers

import (
	"net/http"
	"scifi-search/app/utils"
	"scifi-search/app/views"
)

func registerHistoryHandlers() {

	http.HandleFunc("/history", handleHistory)

}

func handleHistory(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		showHistory(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}

}

func showHistory(w http.ResponseWriter, r *http.Request) {
	user := getCurrentUser(w, r)

	// Si hay un usuario autenticado (con cookies de sesión)...
	if user != nil {

		rows, err := queries.ListHistoricSearchesFromUser(r.Context(), user.UserID)
		if err != nil {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		searches := make([]utils.HistoricSearchView, len(rows))
		for i, r := range rows {
			searches[i] = utils.HistoricSearchView{
				Query: r.SearchString,
				Date:  r.SearchDatetime,
			}
		}

		component := views.HistoryPage(searches, utils.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)

	} else {
		component := views.UnloggedPage(utils.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)
	}
}
