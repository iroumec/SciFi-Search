package languages

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import "golang.org/x/text/language"

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

var SupportedLanguages = map[string]language.Tag{
	"es": language.Spanish,
	"en": language.English,
	// If a new language wants to be added, just add it here,
	// in the array below and create a json file for it.
}

var LanguagesArray = []string{
	"ES",
	"EN",
}

// ------------------------------------------------------------------------------------------------
