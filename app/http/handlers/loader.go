package handlers

import (
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/http/notifications/cookies"
	"scifi-search/app/languages"
	"scifi-search/app/utils/checkers"
	"scifi-search/app/views"
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

	component := views.LoaderPage(auth.GetCurrentAuthorizationLevel(w, r), languages.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
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
	if checkers.IsThereAnEmptyField(name, surname, email, password) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	_, err := createUser(name, surname, email, password, auth.LoaderRole)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := languages.GetTranslatorFromRequest(r)("Loader creado con éxito")
	cookies.AddFlashCookie(w, message)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------
