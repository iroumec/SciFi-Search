package extractors

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"html/template"
	"regexp"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Given an HTML text, it extracts the content between <body> and </body>.
func ExtractBodyContent(html string) template.HTML {
	re := regexp.MustCompile(`(?s)<body[^>]*>(.*?)</body>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) > 1 {
		return template.HTML(matches[1])
	}
	return template.HTML(html)
}

// ------------------------------------------------------------------------------------------------
