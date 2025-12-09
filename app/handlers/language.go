package handlers

import (
	"net/http"

	"scifi-search/app/utils"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func registerLanguageHandlers() {
	http.HandleFunc("/language", setLanguageHandler)
}

func loadAllMessages() error {
	es, err := utils.LoadJSON("locales/es.json")
	if err != nil {
		return err
	}

	en, err := utils.LoadJSON("locales/en.json")
	if err != nil {
		return err
	}

	for k, v := range es {
		message.SetString(language.Spanish, k, v)
	}
	for k, v := range en {
		message.SetString(language.English, k, v)
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	matcher := language.NewMatcher([]language.Tag{
		language.English, // Primer idioma en la lista, por lo que es el idioma por defecto.
		language.Spanish,
	})

	tag, _, _ := language.ParseAcceptLanguage(lang)
	matched, _, _ := matcher.Match(tag...)

	p := message.NewPrinter(matched)
	p.Fprintf(w, "%s", p.Sprintf("hello_msg"))
}

func setLanguageHandler(w http.ResponseWriter, r *http.Request) {

	// Obtención del lenguaje de la URL:
	language := r.URL.Query().Get("language")

	// Validación del lenguaje.
	if language != "es" && language != "en" {
		http.Error(w, "invalid language", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "lang",
		Value: language,
		Path:  "/",
	})
	w.Write([]byte("Idioma actualizado"))
}

func registerMessage(language language.Tag, key string, msg string) {
	message.SetString(language, key, msg)
}
