package utils

// ------------------------------------------------------------------------------------------------

import (
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ------------------------------------------------------------------------------------------------

type Translator func(key string, args ...any) string

// ------------------------------------------------------------------------------------------------

func detectLanguage(r *http.Request) language.Tag {

	// Primer intento: Utilización de cookie explícita.
	if c, err := r.Cookie("language"); err == nil {
		if tag, err := language.Parse(c.Value); err == nil {
			return tag
		}
	}

	// Segundo intento: Utilización de "Accept-Language".
	if al := r.Header.Get("Accept-Language"); al != "" {
		tags, _, _ := language.ParseAcceptLanguage(al)
		if len(tags) > 0 {
			return tags[0]
		}
	}

	// Tercer intento: Fallback a default.
	return language.English
}

// ------------------------------------------------------------------------------------------------

func GetTranslatorFromRequest(r *http.Request) Translator {
	tag := detectLanguage(r)
	p := message.NewPrinter(tag)
	return func(key string, args ...any) string {
		return p.Sprintf(key, args...)
	}
}

// ------------------------------------------------------------------------------------------------
