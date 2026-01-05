package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	sqlc "scifi-search/app/database"
	"scifi-search/app/http/cookies"
	"scifi-search/app/languages"
	"scifi-search/app/utils/structures"
	"scifi-search/app/views"
	"strings"
)

// ------------------------------------------------------------------------------------------------
// Registro de endpoints
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
// Definición de handlers
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
		views.EditPasswordSection(languages.GetTranslatorFromRequest(r)).Render(r.Context(), w)
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
		views.AvatarOptions(languages.GetTranslatorFromRequest(r)).Render(r.Context(), w)
	default:
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
	}

}

// ------------------------------------------------------------------------------------------------

func handleSettingsCancel(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		currentUser, err := getCurrentUser(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		preferences, err := queries.ListPreferencesFromUser(r.Context(), currentUser.UserID)
		if err != nil {
			log.Println("ListPreferencesFromUser:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		user := structures.User{
			Name:            currentUser.Name,
			Surname:         currentUser.Surname,
			AvatarURLString: currentUser.AvatarUrl.String,
			AvatarURLValid:  currentUser.AvatarUrl.Valid,
			Email:           *getCurrentUserEmail(w, r),
		}

		views.SettingsForm(user, preferences, languages.GetTranslatorFromRequest(r)).Render(r.Context(), w)
	default:
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
	}

}

// ------------------------------------------------------------------------------------------------
// Definición de funciones
// ------------------------------------------------------------------------------------------------

// Muestra la página de configuración.
func showSettings(w http.ResponseWriter, r *http.Request) {

	currentUser, err := getCurrentUser(w, r)
	if err != nil {
		component := views.UnloggedPage(languages.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)
	} else {

		preferences, err := queries.ListPreferencesFromUser(r.Context(), currentUser.UserID)
		if err != nil {
			log.Println("ListPreferencesFromUser:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		user := structures.User{
			Name:            currentUser.Name,
			Surname:         currentUser.Surname,
			AvatarURLString: currentUser.AvatarUrl.String,
			AvatarURLValid:  currentUser.AvatarUrl.Valid,
			Email:           *getCurrentUserEmail(w, r),
		}

		component := views.SettingsPage(user, preferences, languages.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)
	}
}

// ------------------------------------------------------------------------------------------------

// Renderiza el campo que está siendo editado.
func editField(w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	//Obtengo el usuario
	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Obtencion de los valores guardados en nuestra BD que cambian
	name := user.Name
	if v := r.Form.Get("name"); v != "" {
		name = v
	}

	surname := user.Surname
	if v := r.Form.Get("surname"); v != "" {
		surname = v
	}

	saveAvatar(user, w, r)
	if r.Form.Get("delete-avatar") == "true" { //En caso de Upload + Delete, gana Delete
		err := deleteAvatar(r.Context(), user.UserID)
		if err != nil {
			http.Error(w, "Error deleting avatar", http.StatusInternalServerError)
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
			cookies.AddFlashCookie(w, "Email actualizado. Por favor, verfiique su nuevo email.")
		} else {
			cookies.AddFlashCookie(w, "Error interno del servidor.")
			return
		}
	}

	// Actualización de la contrasñea.
	currentPassword := r.Form.Get("current-password")
	newPassword := r.Form.Get("new-password")
	if currentPassword != "" && newPassword != "" {
		err := updatePassword(user, currentPassword, newPassword)
		if err != nil {
			cookies.AddFlashCookie(w, "Error interno del servidor.")
			return
		}
	}

	//Renderización final
	currentUpdatedUser, err := getCurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedUser := structures.User{
		Name:            currentUpdatedUser.Name,
		Surname:         currentUpdatedUser.Surname,
		AvatarURLString: currentUpdatedUser.AvatarUrl.String,
		AvatarURLValid:  currentUpdatedUser.AvatarUrl.Valid,
		Email:           *getCurrentUserEmail(w, r),
	}

	component := views.SettingsForm(updatedUser, updatedPreferences, languages.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func saveAvatar(user *sqlc.User, w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "No se pudo leer el archivo", http.StatusBadRequest)
		return
	} else if err == nil {

		defer file.Close()

		// Se cambia el tamaño de la imagen.
		resizedFile, err := ResizeImageToAvatar(file)
		if err != nil {
			http.Error(w, "Error procesando imagen", http.StatusInternalServerError)
			return
		}
		// Se sube el archivo al almacenamiento de objetos.
		url, err := UploadAvatar(r.Context(), bucketName, user.UserID, resizedFile)
		if err != nil {
			http.Error(w, "Error subiendo avatar", http.StatusInternalServerError)
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

// Actualiza las preferencias.
func updatePreferences(user *sqlc.User, r *http.Request) []string {

	if debug {
		log.Println(r.Form["preferences[]"])
	}
	queries.RemoveAllPreferenceFromUser(r.Context(), user.UserID)
	var updatedPreferences []string
	if preferences := r.Form["preferences[]"]; len(preferences) != 0 {
		for _, p := range preferences {
			if p != "" && !structures.Exists(updatedPreferences, p) {
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
		views.InfoField("password", "password", "********", languages.GetTranslatorFromRequest(r)).Render(r.Context(), w)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------
