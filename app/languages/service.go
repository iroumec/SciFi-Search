package languages

import (
	"sort"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

func LoadAllMessagesFromFiles() error {
	for code, tag := range SupportedLanguages {
		messages, err := loadLanguageFile(code)
		if err != nil {
			return err
		}

		registerMessages(tag, messages)
	}
	return nil
}

// ------------------------------------------------------------------------------------------------

func LoadAllMessagesFromFolders() error {
	for code, tag := range SupportedLanguages {
		messages, err := loadLanguageFolder(code)
		if err != nil {
			return err
		}

		registerMessages(tag, messages)
	}
	return nil
}

// ------------------------------------------------------------------------------------------------

func GetSupportedLanguageCodes() []string {
	langs := make([]string, 0, len(SupportedLanguages))
	for k := range SupportedLanguages {
		langs = append(langs, k)
	}
	sort.Strings(langs) // Stable ordering.
	return langs
}

// ------------------------------------------------------------------------------------------------
