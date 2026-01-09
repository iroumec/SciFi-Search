package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"scifi-search/app/auth"
	"scifi-search/app/avatars"
	sqlc "scifi-search/app/database"
	"scifi-search/app/http/notifications/cookies"
	"scifi-search/app/http/notifications/triggers"
	"scifi-search/app/languages"
	"scifi-search/app/utils/getters"
	"scifi-search/app/utils/structures"
	"scifi-search/app/views"
)

// ------------------------------------------------------------------------------------------------
// Registro de endpoints
// ------------------------------------------------------------------------------------------------

func RegisterSettingsHandlers() {

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
			Email:           *auth.GetCurrentUserEmail(w, r),
		}

		views.SettingsForm(user, preferences, auth.GetAuthenticationLevel(currentUser.AuthID), languages.GetTranslatorFromRequest(r)).Render(r.Context(), w)
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
		component := views.UnloggedPage(auth.NoRole.Level, languages.GetTranslatorFromRequest(r))
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
			Email:           *auth.GetCurrentUserEmail(w, r),
		}

		component := views.SettingsPage(user, preferences, auth.GetAuthenticationLevel(currentUser.AuthID), languages.GetTranslatorFromRequest(r))
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

	email := auth.GetCurrentUserEmail(w, r)
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

func saveSettings(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emailUpdated, err := updateUserData(r.Context(), user, w, r)
	if err != nil {
		log.Println(err)
		return
	}

	renderUpdatedSettings(w, r, emailUpdated)
}

// ------------------------------------------------------------------------------------------------

func updateUserData(ctx context.Context, user *sqlc.User, w http.ResponseWriter, r *http.Request) (bool, error) {
	// Actualizar nombre y apellido
	name := getters.GetOrDefault(r.Form.Get("name"), user.Name)
	surname := getters.GetOrDefault(r.Form.Get("surname"), user.Surname)

	queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		UserID:  user.UserID,
		Name:    name,
		Surname: surname,
	})

	// Actualización de avatar.
	saveAvatar(user, w, r)
	if r.Form.Get("delete-avatar") == "true" {
		if err := avatarService.Delete(ctx, user.UserID); err != nil {
			http.Error(w, "Error deleting avatar", http.StatusInternalServerError)
			return false, err
		}
	}

	// Actualización de preferencias.
	updatePreferences(user, r)

	// Actualización de email.
	emailUpdated := false
	if newEmail := r.Form.Get("email"); newEmail != "" {
		if err := auth.UpdateEmail(user.AuthID, newEmail); err != nil {
			return false, err
		}
		emailUpdated = true
		auth.SendVerificationEmail(emailService, user.AuthID, newEmail)
	}

	// Actualización de contraseña.
	currentPwd := r.Form.Get("current-password")
	newPwd := r.Form.Get("new-password")
	if currentPwd != "" && newPwd != "" {
		if err := auth.UpdatePassword(user.AuthID, currentPwd, newPwd); err != nil {
			cookies.AddFlashCookie(w, "Error interno del servidor.")
			return false, err
		}
	}

	return emailUpdated, nil
}

// ------------------------------------------------------------------------------------------------

func renderUpdatedSettings(w http.ResponseWriter, r *http.Request, emailUpdated bool) {

	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := languages.GetTranslatorFromRequest(r)
	message := t("Cambios guardados con éxito.")
	if emailUpdated {
		message += " " + t("Verifique su nuevo correo.")
	}

	triggers.AddHXTrigger(w, r, message)

	component := views.SettingsForm(
		structures.User{
			Name:            user.Name,
			Surname:         user.Surname,
			AvatarURLString: user.AvatarUrl.String,
			AvatarURLValid:  user.AvatarUrl.Valid,
			Email:           *auth.GetCurrentUserEmail(w, r),
		},
		updatePreferences(user, r),
		auth.GetAuthenticationLevel(user.AuthID),
		t,
	)
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func saveAvatar(user *sqlc.User, w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "No se pudo leer el archivo", http.StatusBadRequest)
		return
	}

	if err == http.ErrMissingFile {
		return
	}

	defer file.Close()

	resizedFile, err := avatars.ResizeImageToAvatar(file)
	if err != nil {
		http.Error(w, "Error procesando imagen", http.StatusInternalServerError)
		return
	}

	if avatarService == nil {
		http.Error(w, "Avatar service not configured", http.StatusInternalServerError)
		return
	}

	url, err := avatarService.Upload(r.Context(), user.UserID, resizedFile)
	if err != nil {
		http.Error(w, "Error subiendo avatar", http.StatusInternalServerError)
		return
	}

	_ = queries.UploadAvatar(r.Context(), sqlc.UploadAvatarParams{
		UserID: user.UserID,
		AvatarUrl: sql.NullString{
			String: url,
			Valid:  true,
		},
	})
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
