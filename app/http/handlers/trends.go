package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/languages"
	"scifi-search/app/reporting/graphs"
	"scifi-search/app/utils/extractors"
	"scifi-search/app/views"
)

// ------------------------------------------------------------------------------------------------
// Constantes
// ------------------------------------------------------------------------------------------------

const (
	maxResultsShown = 15 // Cantidad máxima de términos que se presentarán en el gráfico.
)

// ------------------------------------------------------------------------------------------------
// Registro de endpoints
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

	switch r.Method {
	case http.MethodGet:
		showTrendingsGraph(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

// Renderiza y muestra el gráfico de las tendencias.
func showTrendingsGraph(w http.ResponseWriter, r *http.Request) {

	// Se obtienen las N búsquedas más relevantes de las últimas 24 horas.
	trendingSearches, err := queries.GetTrendingSearches(r.Context(), maxResultsShown)
	if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Se genera el gráfico en HTML y se almacena en un buffer.
	buffer, err := graphs.GenerateBarChart(trendingSearches)
	if err != nil {
		http.Error(w, "Error al generar gráfico", http.StatusInternalServerError)
		return
	}

	// Se extrae el "body" del contenido HTML.
	// Interesa únicamente el "body" ya que es el contenido del gráfico.
	// El layout y demás se establece independientemente.
	htmlChart := extractors.ExtractBodyContent(buffer.String())

	// Se renderiza la página resultante.
	component := views.TrendsPage(htmlChart, auth.GetCurrentAuthorizationLevel(w, r, queries), languages.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------
