package languages

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ------------------------------------------------------------------------------------------------
// Types
// ------------------------------------------------------------------------------------------------

type Translator func(key string, args ...any) string

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

func GetTranslatorFromRequest(r *http.Request) Translator {
	tag := DetectLanguage(r)
	p := message.NewPrinter(tag)
	return func(key string, args ...any) string {
		return p.Sprintf(key, args...)
	}
}

// ------------------------------------------------------------------------------------------------
// Functions
// ------------------------------------------------------------------------------------------------

func DetectLanguage(r *http.Request) language.Tag {

	// If the request es nil...
	if r == nil {
		return language.English
	}

	// First try: explicit cookie.
	if c, err := r.Cookie("language"); err == nil {
		if tag, err := language.Parse(c.Value); err == nil {
			return tag
		}
	}

	// Second try: use of "Accept-Language".
	if al := r.Header.Get("Accept-Language"); al != "" {
		tags, _, _ := language.ParseAcceptLanguage(al)
		if len(tags) > 0 {
			return tags[0]
		}
	}

	// Third try: fallback.
	return language.English
}

// ------------------------------------------------------------------------------------------------
