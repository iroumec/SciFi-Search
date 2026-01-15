package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/http/notifications/cookies"
	"scifi-search/app/languages"
	"scifi-search/app/utils/checkers"
	"scifi-search/app/views"
)

// ------------------------------------------------------------------------------------------------
// Handlers
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
// Functions
// ------------------------------------------------------------------------------------------------

func showNewLoaderPage(w http.ResponseWriter, r *http.Request) {

	component := views.LoaderPage(
		auth.GetCurrentAuthorizationLevel(w, r),
		languages.GetTranslatorFromRequest(r),
	)
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func createNewLoader(w http.ResponseWriter, r *http.Request) {

	// Form parsing.
	if err := r.ParseForm(); err != nil {
		http.Error(w, FormParsingError.Error(), http.StatusBadRequest)
		return
	}

	// User data.
	name := r.Form.Get("name")
	surname := r.Form.Get("surname")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Validation.
	if checkers.IsThereAnEmptyField(name, surname, email, password) {
		http.Error(w, RequiredDataNotSpecified.Error(), http.StatusBadRequest)
		return
	}

	_, err := createUser(
		name,
		surname,
		email,
		password,
		auth.LoaderRole,
		languages.GetTranslatorFromRequest(r),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := languages.GetTranslatorFromRequest(r)("new-loader-created")
	cookies.AddFlashCookie(w, message)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------
