package loaders

// ------------------------------------------------------------------------------------------------

import (
	"encoding/json"
	"os"
)

// ------------------------------------------------------------------------------------------------

func LoadJSON(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ------------------------------------------------------------------------------------------------
