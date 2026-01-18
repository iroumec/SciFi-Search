package converters

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"unicode"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Converts the first letter of a string to uppercase.
func UppercaseFirstLetter(s string) string {
	if s == "" {
		return ""
	}
	// Convert string to a slice of runes to handle Unicode characters correctly.
	runes := []rune(s)

	// Convert the first rune to uppercase.
	runes[0] = unicode.ToUpper(runes[0])

	// Convert the slice of runes back to a string.
	return string(runes)
}

// ------------------------------------------------------------------------------------------------
