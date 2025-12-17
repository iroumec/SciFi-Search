package extractors

import (
	"html/template"
	"regexp"
)

// ------------------------------------------------------------------------------------------------

// Se extrae solo el contenido entre <body> y </body>.
// Esto debido a que el Layout ya se encarga de todo lo demás.
func ExtractBodyContent(html string) template.HTML {
	re := regexp.MustCompile(`(?s)<body[^>]*>(.*?)</body>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return template.HTML(matches[1])
	}
	return template.HTML(html)
}

// ------------------------------------------------------------------------------------------------
