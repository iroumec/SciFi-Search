package languages

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"fmt"
	"os"
	"scifi-search/app/utils/loaders"
	"strings"
)

// ------------------------------------------------------------------------------------------------
// Constants
// ------------------------------------------------------------------------------------------------

const languagesDirectory = "resources/languages/"

// ------------------------------------------------------------------------------------------------
// Functions
// ------------------------------------------------------------------------------------------------

func loadLanguageFolder(code string) (map[string]string, error) {
	basePath := fmt.Sprintf("%s%s/", languagesDirectory, code)

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filename := strings.TrimSuffix(entry.Name(), ".json")
		path := basePath + entry.Name()

		messages, err := loaders.LoadJSON(path)
		if err != nil {
			return nil, err
		}

		for key, value := range messages {
			result[filename+"."+key] = value
		}
	}

	return result, nil
}

// ------------------------------------------------------------------------------------------------

func loadLanguageFile(code string) (map[string]string, error) {
	path := fmt.Sprintf("%s%s.json", languagesDirectory, code)
	return loaders.LoadJSON(path)
}

// ------------------------------------------------------------------------------------------------
