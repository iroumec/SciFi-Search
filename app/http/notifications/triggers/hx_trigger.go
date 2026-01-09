package triggers

import "net/http"

func AddHXTrigger(w http.ResponseWriter, r *http.Request, message string) bool {

	// Si es petición HTMX, se envia un trigger para mostrar popup.
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Trigger", `{
					"showFlash": {
						"message": "`+message+`"
					}
				}`)
		return true
	}

	return false
}
