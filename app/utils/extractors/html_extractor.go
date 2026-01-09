package extractors

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"html/template"
	"regexp"
)

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

// Extrae el contenido entre <body> y </body> de un HTML.
func ExtractBodyContent(html string) template.HTML {
	re := regexp.MustCompile(`(?s)<body[^>]*>(.*?)</body>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return template.HTML(matches[1])
	}
	return template.HTML(html)
}

// ------------------------------------------------------------------------------------------------
