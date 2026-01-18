package charts

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"bytes"
	"scifi-search/app/database"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

// ------------------------------------------------------------------------------------------------
// Interfaces
// ------------------------------------------------------------------------------------------------

// Contract for any type of chart.
type ChartRenderer interface {
	Render(xValues []string, yValues []int, seriesName, yAxisName string) (bytes.Buffer, error)
}

// ------------------------------------------------------------------------------------------------
// Functions
// ------------------------------------------------------------------------------------------------

func generateGraph(
	renderer ChartRenderer, values []database.GetTrendingSearchesRow,
	seriesName, yAxisName string,
) (bytes.Buffer, error) {

	xValues := make([]string, len(values))
	yValues := make([]int, len(values))
	for i, row := range values {
		xValues[i] = row.SearchString
		yValues[i] = int(row.Count)
	}

	return renderer.Render(xValues, yValues, seriesName, yAxisName)
}

// ------------------------------------------------------------------------------------------------

// Global options, shared by all charts.
func getGlobalOptions(yAxisName string) []charts.GlobalOpts {
	return []charts.GlobalOpts{
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "400px",
			Theme:  types.ThemeWesteros,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			//Name:         "Términos de búsqueda", // It overlaps with the x-axis names.
			NameLocation: "middle",
			NameGap:      30,
			AxisLabel: &opts.AxisLabel{
				Show:     opts.Bool(true),
				Interval: "0", // All tags are shown.
				Rotate:   -45, // Tags are rotated so they don't overlap.
				FontSize: 12,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:         yAxisName,
			NameLocation: "middle",
			NameGap:      50,
			MinInterval:  1,
			AxisLabel: &opts.AxisLabel{
				Show:     opts.Bool(true),
				FontSize: 12,
			},
			// Manejo de la cuadrícula.
			SplitLine: &opts.SplitLine{
				Show: opts.Bool(true),
				LineStyle: &opts.LineStyle{
					Color: "#333333",
					Width: 1,
					Type:  "solid", // Another option: "dashed".
				},
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    opts.Bool(true),
			Trigger: "axis",
		}),
		charts.WithGridOpts(opts.Grid{
			Left:         "10%",
			Right:        "5%",
			Bottom:       "15%", // More space below the rotated tags.
			Top:          "10%",
			ContainLabel: opts.Bool(true),
		}),
	}
}

// ------------------------------------------------------------------------------------------------
