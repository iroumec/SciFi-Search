package handlers

// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/utils"
	"scifi-search/app/utils/extractors"
	"scifi-search/app/utils/graphs"
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
	trendingSearches, err := queries.GetTrendingSearches(r.Context())
	if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Se limitan los resultados a los N más relevantes.
	limit := maxResultsShown
	if len(trendingSearches) > limit {
		trendingSearches = trendingSearches[:limit]
	}

	buffer, err := graphs.GenerateBarChart(trendingSearches)
	if err != nil {
		http.Error(w, "Error al generar gráfico", http.StatusInternalServerError)
		return
	}

	htmlChart := extractors.ExtractBodyContent(buffer.String())

	views.TrendsPage(htmlChart, isUserAuthenticated(w, r), utils.GetTranslatorFromRequest(r)).Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------
