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

// Bar chart implementations.
type BarChartRenderer struct{}

// ------------------------------------------------------------------------------------------------

func (b *BarChartRenderer) Render(
	xValues []string, yValues []int, seriesName, yAxisName string,
) (bytes.Buffer, error) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(getGlobalOptions(yAxisName)...)
	bar.SetXAxis(xValues).AddSeries(seriesName, b.generateItems(yValues))

	var buffer bytes.Buffer
	err := bar.Render(&buffer)
	return buffer, err
}

// ------------------------------------------------------------------------------------------------

func (b *BarChartRenderer) generateItems(values []int) []opts.BarData {
	items := make([]opts.BarData, len(values))
	for i, v := range values {
		items[i] = opts.BarData{Value: v}
	}
	return items
}

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

func GenerateBarChart(
	values []database.GetTrendingSearchesRow, seriesName, yAxisName string,
) (bytes.Buffer, error) {

	return generateGraph(&BarChartRenderer{}, values, seriesName, yAxisName)
}

// ------------------------------------------------------------------------------------------------
