package graphs

import (
	"bytes"
	"scifi-search/app/database"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// ------------------------------------------------------------------------------------------------

// Implementación del gráfico de barras.
type BarChartRenderer struct{}

func GenerateBarChart(values []database.GetTrendingSearchesRow) (bytes.Buffer, error) {

	return generateGraph(&BarChartRenderer{}, values)
}

func (b *BarChartRenderer) Render(xValues []string, yValues []int) (bytes.Buffer, error) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(getGlobalOptions()...)
	bar.SetXAxis(xValues).AddSeries("Búsquedas", b.generateItems(yValues))

	var buffer bytes.Buffer
	err := bar.Render(&buffer)
	return buffer, err
}

func (b *BarChartRenderer) generateItems(values []int) []opts.BarData {
	items := make([]opts.BarData, len(values))
	for i, v := range values {
		items[i] = opts.BarData{Value: v}
	}
	return items
}
