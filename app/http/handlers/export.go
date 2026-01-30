package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/database"
	"scifi-search/app/export"
	"scifi-search/app/http/notifications"
	"scifi-search/app/languages"
	"scifi-search/app/views"
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Endpoints.
func RegisterExportationHandlers() {

	http.HandleFunc("/documents/export", exportDocumentsHandler)
	http.HandleFunc("/documents/open-export-modal", openExportModal)
}

// ------------------------------------------------------------------------------------------------
// Handlers
// ------------------------------------------------------------------------------------------------

func exportDocumentsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, MethodNotAllowedError.Error(), http.StatusBadRequest)
		return
	}

	exportDocuments(w, r)
}

// ------------------------------------------------------------------------------------------------
// Functions.
// ------------------------------------------------------------------------------------------------

func exportDocuments(w http.ResponseWriter, r *http.Request) {

	format := r.URL.Query().Get("format")

	documents, err := getUserDocuments(w, r)
	if err != nil {
		log.Printf("An error ocurred while obtaining the documents to export.")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the export model was opened, there must be documents.
	if len(documents) <= 0 {
		addEmptyDocumentsNotification(w, r)
		return
	}

	var exporter export.Exporter[database.Document]

	switch format {
	case "csv":
		exporter = export.CSVExporter[database.Document]{}
	case "xlsx", "excel":
		exporter = export.ExcelExporter[database.Document]{}
	default:
		http.Error(w, UnsupportedExportFormatError.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", exporter.ContentType())
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, exporter.FileName()),
	)

	if err := exporter.Export(w, documents); err != nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("An exportation of the documents to %s was executed.", format)
}

// ------------------------------------------------------------------------------------------------

func openExportModal(w http.ResponseWriter, r *http.Request) {

	documents, err := getUserDocuments(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(documents) <= 0 {
		addEmptyDocumentsNotification(w, r)
		return
	}

	views.ExportModal(languages.GetTranslatorFromRequest(r)).Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func getUserDocuments(w http.ResponseWriter, r *http.Request) ([]database.Document, error) {

	currentUser, err := getCurrentUser(w, r)
	if err != nil {
		return nil, err
	}

	var documents []database.Document
	if auth.GetAuthenticationLevel(currentUser.AuthID) >= auth.AdminRole.Level {
		documents, err = queries.ListAllDocuments(r.Context())
	} else {
		documents, err = queries.ListDocumentsByUser(r.Context(),
			sql.NullInt32{Int32: currentUser.UserID, Valid: true})
	}

	if err != nil {
		return nil, InternalServerError
	}

	return documents, nil
}

// ------------------------------------------------------------------------------------------------

func addEmptyDocumentsNotification(w http.ResponseWriter, r *http.Request) {
	translator := languages.GetTranslatorFromRequest(r)
	notifications.ShowFlash(w, r, translator("messages.no-documents-to-export"))
}

// ------------------------------------------------------------------------------------------------
