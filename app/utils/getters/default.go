package getters

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Returns the value if it's not empty. If it's, returns the default value.
func GetOrDefault(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

// ------------------------------------------------------------------------------------------------
