package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/database"
	"scifi-search/app/http/middlewares"
	"scifi-search/app/http/notifications"
	"scifi-search/app/http/notifications/cookies"
	"strconv"
	"strings"

	"scifi-search/app/languages"
	"scifi-search/app/views"

	sqlc "scifi-search/app/database"
)

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

var (
	InvalidIDError      = errors.New("Invalid ID")
	UnknownError        = errors.New("Unknown error")
	InternalServerError = errors.New("Internal server error")
	UserNotFoundError   = errors.New("User not found")
)

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

// Registro de endpoints.
func RegisterFundingHandlers() {

	http.HandleFunc(
		"/funding",
		middlewares.RequiresEmailVerified(
			middlewares.RequiresAuthorization(
				addFundingHandler,
				auth.LoaderRole.Level,
			),
		),
	)

	http.HandleFunc(
		"/funding/open-add-new",
		middlewares.RequiresEmailVerified(
			middlewares.RequiresAuthorization(
				openAddNewFundingModal,
				auth.LoaderRole.Level,
			),
		),
	)

	http.HandleFunc(
		"/funding/open-edit/",
		middlewares.RequiresEmailVerified(
			middlewares.RequiresAuthorization(
				openEditFundingModal,
				auth.LoaderRole.Level,
			),
		),
	)

	http.HandleFunc(
		"/funding/",
		middlewares.RequiresEmailVerified(
			middlewares.RequiresAuthorization(
				addFundingHandler,
				auth.LoaderRole.Level,
			),
		),
	)

	http.HandleFunc(
		"/funding/update-items",
		middlewares.RequiresEmailVerified(
			middlewares.RequiresAuthorization(
				updateFundingList,
				auth.LoaderRole.Level,
			),
		),
	)
}

// ------------------------------------------------------------------------------------------------
// Handlers
// ------------------------------------------------------------------------------------------------

func addFundingHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		showFundingsManagementPage(w, r)
	case http.MethodPost:
		addFunding(w, r)
	case http.MethodDelete:
		deleteFunding(w, r)
	case http.MethodPut:
		editFundingModal(w, r)
	default:
		http.Error(w, MethodNotAllowedError.Error(), http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func showFundingsManagementPage(w http.ResponseWriter, r *http.Request) {

	fundings, totalFundings, err := getFundingDocs(w, r, 1)
	if err != nil {
		return
	}

	component := views.ManageFundingPage(fundings, totalFundings, auth.GetCurrentAuthorizationLevel(w, r), documentTypes, documentAreas, documentCountriesBasedOn, documentGrantors, documentCurrencies, languages.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func getFundingDocs(w http.ResponseWriter, r *http.Request, offset int) ([]map[string]any, int, error) {

	var fundings []map[string]any
	var fundingsDocs []database.Document
	var totalFundings int64

	user, err := getCurrentUser(w, r)
	if err != nil {
		cookies.AddFlashCookie(w, languages.GetTranslatorFromRequest(r)("Ha ocurrido un error inesperado."))
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, 0, err
	}

	if auth.GetAuthenticationLevel(user.AuthID) == auth.AdminRole.Level {

		// De ser admin, se listan todos los documentos.
		fundingsDocs, err = queries.ListAllDocuments(r.Context(), sqlc.ListAllDocumentsParams{
			Limit:  10,
			Offset: int32((offset - 1) * 10),
		})

		totalFundings, err = queries.CountAllDocuments(r.Context())
		if err != nil {
			log.Fatal("Error al contar los documentos")
		}

	} else {

		// De ser loader, se listan solo sus documentos.
		fundingsDocs, err = queries.ListDocumentsByUser(r.Context(), sqlc.ListDocumentsByUserParams{
			UserID: sql.NullInt32{Int32: user.UserID, Valid: true},
			Limit:  10,
			Offset: int32((offset - 1) * 10),
		})

		totalFundings, err = queries.CountDocumentsByUser(r.Context(), sql.NullInt32{Int32: user.UserID, Valid: true})
		if err != nil {
			log.Fatal("Error al contar los documentos")
		}
	}
	if err != nil {
		return nil, 0, err
	}

	for _, doc := range fundingsDocs {
		funding := map[string]any{
			"id":          doc.ID,
			"Nombre":      doc.Name,
			"Tipo":        doc.Type,
			"Gran area 1": doc.FirstArea,
			"Gran area 2": doc.SecondArea.String,
			"Link":        doc.Link.String,
			"Descripcion": doc.Description.String,
			"Pais":        doc.BasedOn.String,
			"Otorgante":   doc.Grantor.String,
			"Moneda":      doc.Currency,
			"Monto":       doc.Amount,
			"Deadline":    doc.Deadline,
		}
		fundings = append(fundings, funding)
	}

	return fundings, int(totalFundings), nil
}

// ------------------------------------------------------------------------------------------------

func updateFundingList(w http.ResponseWriter, r *http.Request) {

	page := getPage(r)

	fundings, totalFundings, err := getFundingDocs(w, r, page)
	if err != nil {
		return
	}

	component := views.FundingList(fundings, page, totalFundings, languages.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func getPage(r *http.Request) int {
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1 // TODO: dejarlo así o levantamos error?
	}

	return page
}

// ------------------------------------------------------------------------------------------------

func addFunding(w http.ResponseWriter, r *http.Request) {

	// Parsing del formulario.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Obtención de los datos del formulario.
	name := r.Form.Get("name")
	fundingType := r.Form.Get("type")
	firstArea := r.Form.Get("first-area")
	secondArea := r.Form.Get("second-area")
	link := r.Form.Get("link")
	description := r.Form.Get("description")
	basedOn := r.Form.Get("based-on")
	grantor := r.Form.Get("grantor")
	currency := r.Form.Get("currency")
	amount := r.Form.Get("amount")
	deadline := r.Form.Get("deadline")

	user, err := getCurrentUser(w, r)
	if err != nil {
		cookies.AddFlashCookie(w, languages.GetTranslatorFromRequest(r)("Ha ocurrido un error inesperado."))
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	document, err := queries.AddDocument(r.Context(), sqlc.AddDocumentParams{
		Name:        name,
		UserID:      sql.NullInt32{Int32: user.UserID, Valid: true},
		Type:        fundingType,
		FirstArea:   firstArea,
		SecondArea:  sql.NullString{String: secondArea, Valid: secondArea != ""},
		Link:        sql.NullString{String: link, Valid: link != ""},
		Description: sql.NullString{String: description, Valid: description != ""},
		BasedOn:     sql.NullString{String: basedOn, Valid: basedOn != ""},
		Grantor:     sql.NullString{String: grantor, Valid: grantor != ""},
		Currency:    currency,
		Amount:      amount,
		Deadline:    deadline,
	})

	_, err = client.Index(indexName).AddDocuments(map[string]any{
		"id":          document.ID,
		"Usuario":     document.UserID,
		"Nombre":      document.Name,
		"Tipo":        document.Type,
		"Gran area 1": document.FirstArea,
		"Gran area 2": document.SecondArea.String,
		"Link":        document.Link.String,
		"Descripcion": document.Description.String,
		"Pais":        document.BasedOn.String,
		"Otorgante":   document.Grantor.String,
		"Moneda":      document.Currency,
		"Monto":       document.Amount,
		"Deadline":    document.Deadline,
	}, &primaryKey)
	if err != nil {
		log.Fatal(err)
	}

	notifyFundingAddition(r.Context(), document)

	notifications.ShowFlash(w, r, "New funding added successfully!")

	updateFundingList(w, r)
}

// ------------------------------------------------------------------------------------------------

func notifyFundingAddition(ctx context.Context, document sqlc.Document) {

	users, err := queries.UsersInterestedInFunding(ctx, document.ID)
	if err != nil {
		log.Println("Error obteniendo usuarios interesados:", err)
		return
	}

	emailBody := strings.Join([]string{
		"Se ha publicado un nuevo financiamiento que podría interesarte de acuerdo a tus preferencias: \n\n",
		"Nombre: " + document.Name + "\n\n",
		"Descripción: " + document.Description.String + "\n",
		"\n¡Accede a SciFi para ver sus detalles!\n",
		"\nSi el documento no te interesa, puedes ajustar tus preferencias en la configuración de tu perfil",
	}, "")

	for _, user := range users {

		emailVerified, err := auth.IsEmailVerified(user.AuthID)
		if err != nil || !*emailVerified {
			continue
		}

		userEmail := auth.GetUserEmail(user.AuthID)
		if userEmail == nil {
			continue
		}

		emailService.Send(
			*userEmail,
			"Nuevo financiamiento disponible",
			emailBody,
		)
	}
}

// ------------------------------------------------------------------------------------------------

func deleteFunding(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/funding/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, InvalidIDError.Error(), http.StatusBadRequest)
		return
	}

	err = queries.RemoveDocument(r.Context(), int32(id))
	if err != nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	_, err = client.Index(indexName).DeleteDocument(idStr)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(""))
}

// ------------------------------------------------------------------------------------------------

func openAddNewFundingModal(w http.ResponseWriter, r *http.Request) {
	component := views.FundingModal("new", nil, documentTypes, documentAreas, documentCountriesBasedOn, documentGrantors, documentCurrencies, languages.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func openEditFundingModal(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/funding/open-edit/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, InvalidIDError.Error(), http.StatusBadRequest)
		return
	}

	sqlcDocument, err := queries.GetDocumentByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	document := map[string]any{
		"id":          strconv.Itoa(int(sqlcDocument.ID)),
		"Nombre":      sqlcDocument.Name,
		"Tipo":        sqlcDocument.Type,
		"Gran area 1": sqlcDocument.FirstArea,
		"Gran area 2": sqlcDocument.SecondArea.String,
		"Link":        sqlcDocument.Link.String,
		"Descripcion": sqlcDocument.Description.String,
		"Pais":        sqlcDocument.BasedOn.String,
		"Otorgante":   sqlcDocument.Grantor.String,
		"Moneda":      sqlcDocument.Currency,
		"Monto":       sqlcDocument.Amount,
		"Deadline":    sqlcDocument.Deadline,
	}

	component := views.FundingModal("edit", document, documentTypes, documentAreas, documentCountriesBasedOn, documentGrantors, documentCurrencies, languages.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func editFundingModal(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/funding/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, InvalidIDError.Error(), http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	fundingType := r.Form.Get("type")
	firstArea := r.Form.Get("first-area")
	secondArea := r.Form.Get("second-area")
	link := r.Form.Get("link")
	description := r.Form.Get("description")
	basedOn := r.Form.Get("based-on")
	grantor := r.Form.Get("grantor")
	currency := r.Form.Get("currency")
	amount := r.Form.Get("amount")
	deadline := r.Form.Get("deadline")

	user, err := getCurrentUser(w, r)
	if err != nil {
		cookies.AddFlashCookie(w, languages.GetTranslatorFromRequest(r)("Ha ocurrido un error inesperado."))
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = queries.UpdateDocument(r.Context(), sqlc.UpdateDocumentParams{
		ID:          int32(id),
		Name:        name,
		UserID:      sql.NullInt32{Int32: user.UserID, Valid: true},
		Type:        fundingType,
		FirstArea:   firstArea,
		SecondArea:  sql.NullString{String: secondArea, Valid: secondArea != ""},
		Link:        sql.NullString{String: link, Valid: link != ""},
		Description: sql.NullString{String: description, Valid: description != ""},
		BasedOn:     sql.NullString{String: basedOn, Valid: basedOn != ""},
		Grantor:     sql.NullString{String: grantor, Valid: grantor != ""},
		Currency:    currency,
		Amount:      amount,
		Deadline:    deadline,
	})
	if err != nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	_, err = client.Index(indexName).UpdateDocuments(map[string]any{
		"id":          id,
		"Usuario":     user.UserID,
		"Nombre":      name,
		"Tipo":        fundingType,
		"Gran area 1": firstArea,
		"Gran area 2": secondArea,
		"Link":        link,
		"Descripcion": description,
		"Pais":        basedOn,
		"Otorgante":   grantor,
		"Moneda":      currency,
		"Monto":       amount,
		"Deadline":    deadline,
	}, &primaryKey)
	if err != nil {
		log.Fatal(err)
	}

	notifications.ShowFlash(w, r, "Funding edited successfully!")

	updateFundingList(w, r)
}
