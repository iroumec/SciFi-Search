package graphs

// ------------------------------------------------------------------------------------------------
// Importaciones
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

// Interfaz que define el contrato para cualquier tipo de gráfico
type ChartRenderer interface {
	Render(xValues []string, yValues []int) (bytes.Buffer, error)
}

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

func generateGraph(renderer ChartRenderer, values []database.GetTrendingSearchesRow) (bytes.Buffer, error) {

	xValues := make([]string, len(values))
	yValues := make([]int, len(values))
	for i, row := range values {
		xValues[i] = row.SearchString
		yValues[i] = int(row.Count)
	}

	// Usar el renderer configurado (polimorfismo)
	return renderer.Render(xValues, yValues)
}

// ------------------------------------------------------------------------------------------------

// Opciones globales compartidas por todos los tipos de gráficos.
func getGlobalOptions() []charts.GlobalOpts {
	return []charts.GlobalOpts{
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "100%",
			Height: "400px",
			Theme:  types.ThemeWesteros,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			//Name:         "Términos de búsqueda", // Comentado ya que se solapa con los nombres de las barras.
			NameLocation: "middle",
			NameGap:      30,
			AxisLabel: &opts.AxisLabel{
				Show:     opts.Bool(true),
				Interval: "0", // Se muestran todas las etiquetas.
				Rotate:   -45, // Se rotan las etiquetas para evitar solapamiento.
				FontSize: 12,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:         "Cantidad de búsquedas",
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
					Color: "#333333", // "000000" para negro puro.
					Width: 1,
					Type:  "solid", // "dashed" para líneas punteadas.
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
			Bottom:       "15%", // Más espacio abajo para las etiquetas rotadas.
			Top:          "10%",
			ContainLabel: opts.Bool(true),
		}),
	}
}

// ------------------------------------------------------------------------------------------------
