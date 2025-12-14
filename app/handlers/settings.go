package handlers

import (
	"context"
	"net/http"
	"scifi-search/app/utils"
	"scifi-search/app/views"
	"strings"
	sqlc "scifi-search/app/database"
	//emailpassword "github.com/supertokens/supertokens-golang/recipe/emailpassword"
)

func registerSettingsHandlers() {

	http.HandleFunc("/settings", handleSettings)
	http.HandleFunc("/settings/edit/", handleEditField)
	http.HandleFunc("/settings/save", saveSettings)
	http.HandleFunc("/settings/edit/password", handlePasswordEdit)
	http.HandleFunc("/settings/edit/password/cancel", cancelPasswordEdit)
	http.HandleFunc("/settings/add-preference", handlePreferences)
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

func handlePasswordEdit(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		views.EditPasswordSection().Render(r.Context(), w)
	default: 
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}

}

func handlePreferences(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		r.ParseForm()

		prefs := r.Form["preferences[]"]
		if len(prefs) == 0 {
			http.Error(w, "Missing preference", http.StatusBadRequest)
			return
		}

		views.PreferenceItem(prefs[0]).Render(r.Context(), w)
	default: 
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}

}

func showSettings(w http.ResponseWriter, r *http.Request) {

	user := getCurrentUser(w, r)

	if user != nil {

		preferences, err := queries.ListPreferencesFromUser(context.Background(),user.UserID)
		if err != nil {//TODO: ver si esta bien este rror
			http.Error(w, "Unknown error", http.StatusInternalServerError)
        	return
		}
		component := views.SettingsPage(*user, *getCurrentUserEmail(w, r), preferences, utils.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)

	} else {

		component := views.SettingsPageError()
		component.Render(r.Context(), w)

	}

}

func editField(w http.ResponseWriter, r *http.Request) {

	user := getCurrentUser(w, r)
	email := getCurrentUserEmail(w, r)
	fieldStr := strings.TrimPrefix(r.URL.Path, "/settings/edit/")

	var value,typeStr string
    switch fieldStr {
    case "name":
        value = user.Name
		typeStr = "text"
    case "surname":
        value = user.Surname
		typeStr = "text"
    case "email":
        value = *email
		typeStr = "email"
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

	if v := r.Form.Get("email"); v != "" {
        updateEmail(user, v)
    }

	updatedUser := getCurrentUser(w, r)
	updatedEmail := getCurrentUserEmail(w, r)
	preferences, _ := queries.ListPreferencesFromUser(context.Background(),updatedUser.UserID)

	component := views.SettingsForm(*updatedUser, *updatedEmail, preferences)
	component.Render(r.Context(), w)

}

func updateEmail(user *sqlc.User, newEmail string) {
	
}

func cancelPasswordEdit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		views.InfoField("password", "password", "********").Render(r.Context(),w)
	default: 
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}