package utils

import "net/http"

func HasGETRequestParameters(r *http.Request) bool {

	return len(r.URL.Query()) > 0
}

func HasPOSTRequestParameters(r *http.Request) bool {

	// Se parsea el body de la request.
	err := r.ParseForm()
	if err != nil {
		// Cuerpo mal formado.
		return false
	}

	return len(r.PostForm) > 0
}
