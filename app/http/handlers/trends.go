package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"errors"
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/languages"
	"scifi-search/app/reporting/graphs"
	"scifi-search/app/utils/extractors"
	"scifi-search/app/views"
)

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

// Errors.
var (
	ChartGenerationFailerError = errors.New("error.chart-generation-failed")
)

// ------------------------------------------------------------------------------------------------
// Constants
// ------------------------------------------------------------------------------------------------

const (
	maxResultsShown = 15 // Max ammount of terms showed in the graph.
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Registra los handlers correspondientes a las tendencias.
func RegisterTrendsHandlers() {
	http.HandleFunc("/trends", trendsHandler)
}

// ------------------------------------------------------------------------------------------------
// Handlers
// ------------------------------------------------------------------------------------------------

// Maneja la petición de la página de tendencias.
func trendsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowedError.Error(), http.StatusMethodNotAllowed)
		return
	}

	showTrendingsGraph(w, r)
}

// ------------------------------------------------------------------------------------------------
// Functions
// ------------------------------------------------------------------------------------------------

// Renderiza y muestra el gráfico de las tendencias.
func showTrendingsGraph(w http.ResponseWriter, r *http.Request) {

	// Obtention of the N more relevants searches from the last 24 hours.
	trendingSearches, err := queries.GetTrendingSearches(r.Context(), maxResultsShown)
	if err != nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	// Generation of the HTML graph and storage in a buffer.
	buffer, err := graphs.GenerateBarChart(trendingSearches)
	if err != nil {
		http.Error(w, ChartGenerationFailerError.Error(), http.StatusInternalServerError)
		return
	}

	// Extraction of the body of the HTML content.
	// The body corresponds to the graph content.
	// The layout and other things are established independently.
	htmlChart := extractors.ExtractBodyContent(buffer.String())

	// Rendering of the page.
	component := views.TrendsPage(
		htmlChart,
		auth.GetCurrentAuthorizationLevel(w, r),
		languages.GetTranslatorFromRequest(r),
	)
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------
