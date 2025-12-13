package handlers

import (
	"context"
	"net/http"
	"scifi-search/app/utils"
	"scifi-search/app/views"
	"strings"
	sqlc "scifi-search/app/database"
)

func registerSettingsHandlers() {

	http.HandleFunc("/settings", handleSettings)
	http.HandleFunc("/settings/edit/", handleEditField)
	http.HandleFunc("/settings/save", saveSettings)

}

func handleSettings(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		showSettings(w, r)
	default: 
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}

}

func handleEditField(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		editField(w, r)
	default: 
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}

}

func showSettings(w http.ResponseWriter, r *http.Request) {

	user := getCurrentUser(w, r)

	if user != nil {

		component := views.SettingsPage(*user, utils.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)

	} else {

		component := views.SettingsPageError()
		component.Render(r.Context(), w)

	}

}

func editField(w http.ResponseWriter, r *http.Request) {

	user := getCurrentUser(w, r)
	fieldStr := strings.TrimPrefix(r.URL.Path, "/settings/edit/")

	var value,typeStr string
    switch fieldStr {
    case "name":
        value = user.Name
		typeStr = "text"
    case "surname":
        value = user.Surname
		typeStr = "text"
    /*case "email":
        value = 
		typeStr = "email"
	case "password":
		value =
		typeStr = "password"*/
    default:
        http.Error(w, "invalid field", http.StatusBadRequest)
        return
    }

	component := views.EditableInfoField(fieldStr,typeStr,value)
	component.Render(r.Context(), w)

}

func saveSettings(w http.ResponseWriter, r *http.Request) {

    r.ParseForm()
	user := getCurrentUser(w, r)

	name := user.Name
	if v := r.Form.Get("name"); v != "" {
        name = v
    }

	surname:= user.Surname
    if v := r.Form.Get("surname"); v != "" {
        surname = v
    }

	queries.UpdateUser(context.Background(),sqlc.UpdateUserParams{
		UserID: 	user.UserID,
		Name:		name,
		Surname:	surname,
	})

	updatedUser := getCurrentUser(w, r)

	component := views.SettingsForm(*updatedUser)
	component.Render(r.Context(), w)
}