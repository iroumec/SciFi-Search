package handlers

// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/languages"
	"scifi-search/app/reporting/graphs"
	"scifi-search/app/utils/extractors"
	"scifi-search/app/views"
)

// ------------------------------------------------------------------------------------------------

const (
	maxResultsShown = 15
)

// ------------------------------------------------------------------------------------------------

func registerTrendsHandler() {
	http.HandleFunc("/trends", trendsHandler)
}

// ------------------------------------------------------------------------------------------------

func trendsHandler(w http.ResponseWriter, r *http.Request) {
	trendingSearches, err := queries.GetTrendingSearches(r.Context(), maxResultsShown)
	if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	buffer, err := graphs.GenerateBarChart(trendingSearches)
	if err != nil {
		http.Error(w, "Error al generar gráfico", http.StatusInternalServerError)
		return
	}

	htmlChart := extractors.ExtractBodyContent(buffer.String())

	views.TrendsPage(htmlChart, isUserAuthenticated(w, r), languages.GetTranslatorFromRequest(r)).Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------
