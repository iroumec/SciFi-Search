package handlers

// ------------------------------------------------------------------------------------------------

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	sqlc "scifi-search/app/database"
	"scifi-search/app/utils"
	"scifi-search/app/views"
	"strings"
)

// ------------------------------------------------------------------------------------------------
// Registro de endpoints.
// ------------------------------------------------------------------------------------------------

func registerSettingsHandlers() {

	http.HandleFunc("/settings", handleSettings)
	http.HandleFunc("/settings/edit/", handleEditField)
	http.HandleFunc("/settings/save", saveSettings)
	http.HandleFunc("/settings/edit/password", handlePasswordEdit)
	http.HandleFunc("/settings/edit/password/cancel", cancelPasswordEdit)
	http.HandleFunc("/settings/add-preference", handlePreferences)
	http.HandleFunc("/settings/modify-avatar", handleAvatarModification)
	http.HandleFunc("/settings/cancel", handleSettingsCancel)
}

// ------------------------------------------------------------------------------------------------
// Definición de handlers.
// ------------------------------------------------------------------------------------------------

func handleSettings(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		showSettings(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func handleEditField(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		editField(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func handlePasswordEdit(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		views.EditPasswordSection().Render(r.Context(), w)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

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

// ------------------------------------------------------------------------------------------------

func handleAvatarModification(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		views.AvatarOptions().Render(r.Context(), w)
	default:
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
	}

}

// ------------------------------------------------------------------------------------------------

func handleSettingsCancel(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		user := getCurrentUser(w, r)
		email := getCurrentUserEmail(w, r)
		preferences, err := queries.ListPreferencesFromUser(r.Context(), user.UserID)
		if err != nil {
			log.Println("ListPreferencesFromUser:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		views.SettingsForm(*user, *email, preferences).Render(r.Context(), w)
	default:
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
	}

}

// ------------------------------------------------------------------------------------------------
// Function definitions.
// ------------------------------------------------------------------------------------------------

// Muestra la página de configuración.
func showSettings(w http.ResponseWriter, r *http.Request) {

	user := getCurrentUser(w, r)

	if user != nil {

		preferences, err := queries.ListPreferencesFromUser(r.Context(), user.UserID)
		if err != nil {
			log.Println("ListPreferencesFromUser:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		component := views.SettingsPage(*user, *getCurrentUserEmail(w, r), preferences, utils.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)

	} else {

		component := views.SettingsPageError()
		component.Render(r.Context(), w)

	}
}

// ------------------------------------------------------------------------------------------------

// Renderiza el campo que está siendo editado.
func editField(w http.ResponseWriter, r *http.Request) {

	user := getCurrentUser(w, r)
	email := getCurrentUserEmail(w, r)
	fieldStr := strings.TrimPrefix(r.URL.Path, "/settings/edit/")

	var value, typeStr string
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

	component := views.EditableInfoField(fieldStr, typeStr, value)
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

// Guarda la nueva configuración.
func saveSettings(w http.ResponseWriter, r *http.Request) {

	log.Println(r)
	//Parseo del formulario
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("d2sadsad")
		log.Println(err)
		http.Error(w, "Invalid form", http.StatusBadRequest)

		return
	}
	log.Println("195")
	//Obtengo el usuario
	user := getCurrentUser(w, r)

	//Obtencion de los valores guardados en nuestra BD que cambian
	name := user.Name
	if v := r.Form.Get("name"); v != "" {
		name = v
	}

	surname := user.Surname
	if v := r.Form.Get("surname"); v != "" {
		surname = v
	}
	log.Println("209")
	saveAvatar(user, w, r)
	if r.Form.Get("delete-avatar") == "true" { //En caso de Upload + Delete, gana Delete
		err := deleteAvatar(r.Context(), user.UserID)
		if err != nil {
			http.Error(w, "Error deleting avatar", 500)
			return
		}
	}

	//Actualizacion de datos que se guardan en nuestra BD
	queries.UpdateUser(context.Background(), sqlc.UpdateUserParams{
		UserID:  user.UserID,
		Name:    name,
		Surname: surname,
	})

	updatedPreferences := updatePreferences(user, r)

	// Actualización del email.
	if v := r.Form.Get("email"); v != "" {
		err := updateEmail(user, v)
		if err == nil {
			utils.AddFlashCookie(w, "Email actualizado. Por favor, verfiique su nuevo email.")
		} else {
			utils.AddFlashCookie(w, "Error interno del servidor.")
			return
		}
	}
	log.Println("238")
	// Actualización de la contrasñea.
	currentPassword := r.Form.Get("current-password")
	newPassword := r.Form.Get("new-password")
	if currentPassword != "" && newPassword != "" {
		err := updatePassword(user, currentPassword, newPassword)
		if err != nil {
			utils.AddFlashCookie(w, "Error interno del servidor.")
			return
		}
	}

	//Renderización final
	updatedUser := getCurrentUser(w, r)
	updatedEmail := getCurrentUserEmail(w, r)
	log.Println("253")
	component := views.SettingsForm(*updatedUser, *updatedEmail, updatedPreferences)
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func saveAvatar(user *sqlc.User, w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "No se pudo leer el archivo", 400)
		return
	} else if err == nil {

		defer file.Close()

		// Se cambia el tamaño de la imagen.
		resizedFile, err := ResizeImageToAvatar(file)
		if err != nil {
			http.Error(w, "Error procesando imagen", 500)
			return
		}
		// Se sube el archivo al almacenamiento de objetos.
		url, err := UploadAvatar(r.Context(), bucketName, user.UserID, resizedFile)
		if err != nil {
			http.Error(w, "Error subiendo avatar", 500)
			log.Printf("%s", err)
			return
		}

		// Guardado de la URL en la Base de Datos.
		err = queries.UploadAvatar(r.Context(), sqlc.UploadAvatarParams{
			UserID: user.UserID,
			AvatarUrl: sql.NullString{
				String: url,
				Valid:  true,
			},
		})
	}
}

// ------------------------------------------------------------------------------------------------

func updatePreferences(user *sqlc.User, r *http.Request) []string {

	log.Println(r.Form["preferences[]"])
	queries.RemoveAllPreferenceFromUser(r.Context(), user.UserID)
	var updatedPreferences []string
	if preferences := r.Form["preferences[]"]; len(preferences) != 0 {
		for _, p := range preferences {
			if p != "" && !utils.Exists(updatedPreferences, p) {
				queries.AddPreference(r.Context(), sqlc.AddPreferenceParams{
					UserID:     user.UserID,
					Preference: strings.ToLower(p),
				})
				updatedPreferences = append(updatedPreferences, p)
			}
		}
	}

	return updatedPreferences
}

// ------------------------------------------------------------------------------------------------

// Cancela la edición del campo de la contraseña.
func cancelPasswordEdit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		views.InfoField("password", "password", "********").Render(r.Context(), w)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------
