package languages

// ------------------------------------------------------------------------------------------------
// Service
// ------------------------------------------------------------------------------------------------

func LoadAllMessages() error {
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
