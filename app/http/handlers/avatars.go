package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"database/sql"
	"net/http"

	"scifi-search/app/avatars"
	"scifi-search/app/database"
)

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

// Registro de endpoints.
func RegisterAvatarHandlers() {

	http.HandleFunc("/avatar", uploadAvatar)
}

// ------------------------------------------------------------------------------------------------

func uploadAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Archivo inválido", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Archivo inválido", http.StatusBadRequest)
		return
	}
	defer file.Close()

	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Error(w, "No autenticado", http.StatusUnauthorized)
		return
	}

	data, err := avatars.ResizeImageToAvatar(file)
	if err != nil {
		http.Error(w, "Imagen inválida", http.StatusBadRequest)
		return
	}

	url, err := avatarService.Upload(r.Context(), user.UserID, data)
	if err != nil {
		http.Error(w, "Error subiendo avatar", http.StatusInternalServerError)
		return
	}

	_ = queries.UploadAvatar(r.Context(), database.UploadAvatarParams{
		UserID: user.UserID,
		AvatarUrl: sql.NullString{
			String: url,
			Valid:  true,
		},
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------
