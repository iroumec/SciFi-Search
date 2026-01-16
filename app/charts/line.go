package charts

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"bytes"
	"scifi-search/app/database"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// ------------------------------------------------------------------------------------------------
// Structures
// ------------------------------------------------------------------------------------------------

// Line chart implementation.
type LineChartRenderer struct{}

// ------------------------------------------------------------------------------------------------

func (l *LineChartRenderer) Render(xValues []string, yValues []int) (bytes.Buffer, error) {
	line := charts.NewLine()
	line.SetGlobalOptions(getGlobalOptions()...)
	line.SetXAxis(xValues).AddSeries("Búsquedas", l.generateItems(yValues))

	var buffer bytes.Buffer
	err := line.Render(&buffer)
	return buffer, err
}

// ------------------------------------------------------------------------------------------------

func (l *LineChartRenderer) generateItems(values []int) []opts.LineData {
	items := make([]opts.LineData, len(values))
	for i, v := range values {
		items[i] = opts.LineData{Value: v}
	}
	return items
}

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

func GenerateLineChart(values []database.GetTrendingSearchesRow) (bytes.Buffer, error) {

	return generateGraph(&LineChartRenderer{}, values)
}

// ------------------------------------------------------------------------------------------------
