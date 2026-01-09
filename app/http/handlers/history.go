package handlers

// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/languages"
	"scifi-search/app/utils/structures"
	"scifi-search/app/views"
)

// ------------------------------------------------------------------------------------------------

// Registra los endpoint asociados al historial.
func RegisterHistoryHandlers() {

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
	user, err := getCurrentUser(w, r)
	if err != nil {
		component := views.UnloggedPage(auth.NoRole.Level, languages.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)
	} else {

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

		component := views.HistoryPage(searches, auth.GetAuthenticationLevel(user.AuthID), languages.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)
	}
}

// ------------------------------------------------------------------------------------------------
