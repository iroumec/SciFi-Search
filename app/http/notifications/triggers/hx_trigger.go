package triggers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"strconv"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

func AddHXTrigger(w http.ResponseWriter, r *http.Request, message string) bool {
	if r.Header.Get("HX-Request") != "true" {
		return false
	}

	asciiMessage := strconv.QuoteToASCII(message)
	asciiMessage = asciiMessage[1 : len(asciiMessage)-1] // Quotes are eliminated.

	payload := `{"showFlash":{"message":"` + asciiMessage + `"}}`

	w.Header().Set("HX-Trigger", payload)
	return true
}

// ------------------------------------------------------------------------------------------------
