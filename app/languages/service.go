package languages

// ------------------------------------------------------------------------------------------------
// Service
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
