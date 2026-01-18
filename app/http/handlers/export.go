package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"fmt"
	"log"
	"net/http"
	"scifi-search/app/database"
	"scifi-search/app/export"
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

	documents, err := queries.GetDocuments(r.Context())
	if err != nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
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
	views.ExportModal(languages.GetTranslatorFromRequest(r)).Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------
