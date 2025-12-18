package languages

// ------------------------------------------------------------------------------------------------

import (
	"scifi-search/app/utils/loaders"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ------------------------------------------------------------------------------------------------

const (
	languagesDirectory = "resources/languages/"
)

// ------------------------------------------------------------------------------------------------

func LoadAllMessages() error {
	es, err := loaders.LoadJSON(languagesDirectory + "es.json")
	if err != nil {
		return err
	}

	en, err := loaders.LoadJSON(languagesDirectory + "en.json")
	if err != nil {
		return err
	}

	for key, value := range es {
		message.SetString(language.Spanish, key, value)
	}
	for key, value := range en {
		message.SetString(language.English, key, value)
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
