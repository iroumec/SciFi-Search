package handlers

import (
	"log"
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/utils/checkers"
)

// ------------------------------------------------------------------------------------------------

func loaderHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		showNewLoaderPage(w, r)
	case http.MethodPost:
		createNewLoader(w, r)
	}

}

// ------------------------------------------------------------------------------------------------

func showNewLoaderPage(w http.ResponseWriter, r *http.Request) {

}

// ------------------------------------------------------------------------------------------------

func createNewLoader(w http.ResponseWriter, r *http.Request) {
	// Parseo del formulario enviado por POST.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Obtención de los datos del usuario.
	name := r.Form.Get("name")
	surname := r.Form.Get("surname")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Validación.
	if checkers.HayCampoIncompleto(name, surname, email, password) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	createLoader(name, surname, email, password)
}

// ------------------------------------------------------------------------------------------------

func createLoader(name, surname, email, password string) {

	user, resp := createUser(name, surname, email, password, auth.LoaderRole.Name)

	if user == nil || resp == nil {
		log.Fatal("Ocurrió un error al momento de crear al usuario.")
	}

	log.Printf("Loader created: %s", email)
}

// ------------------------------------------------------------------------------------------------
