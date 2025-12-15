package handlers

// ------------------------------------------------------------------------------------------------

import (
	"bytes"
	"html/template"
	"net/http"
	"regexp"
	"scifi-search/app/database"
	"scifi-search/app/utils"
	"scifi-search/app/views"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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

	buffer, err := generateGraph(trendingSearches)
	if err != nil {
		http.Error(w, "Error al generar gráfico", 500)
		return
	}

	htmlChart := extractBodyContent(buffer.String())

	views.TrendsPage(htmlChart, isUserAuthenticated(w, r), utils.GetTranslatorFromRequest(r)).Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func generateGraph(trendingSearches []database.GetTrendingSearchesRow) (bytes.Buffer, error) {

	xValues := make([]string, len(trendingSearches))
	yValues := make([]int, len(trendingSearches))
	for i, row := range trendingSearches {
		xValues[i] = row.SearchString
		yValues[i] = int(row.Count)
	}

	line := charts.NewLine()
	line.SetGlobalOptions(
		//charts.WithTitleOpts(opts.Title{Title: "Tendencias de búsquedas"}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "400px",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			MinInterval: 1, // Esto fuerza intervalos mínimos de 1. Es decir, no se muestran decimales.
		}),
	)
	line.SetXAxis(xValues).
		AddSeries("Búsquedas", generateLineItems(yValues))

	var buffer bytes.Buffer

	err := line.Render(&buffer)

	return buffer, err
}

// ------------------------------------------------------------------------------------------------

func generateLineItems(values []int) []opts.LineData {
	items := make([]opts.LineData, len(values))
	for i, v := range values {
		items[i] = opts.LineData{Value: v}
	}
	return items
}

// ------------------------------------------------------------------------------------------------

// Se extrae solo el contenido entre <body> y </body>.
// Esto debido a que el Layout ya se encarga de todo lo demás.
func extractBodyContent(html string) template.HTML {
	re := regexp.MustCompile(`(?s)<body[^>]*>(.*?)</body>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return template.HTML(matches[1])
	}
	return template.HTML(html)
}

// ------------------------------------------------------------------------------------------------
