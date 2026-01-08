package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"scifi-search/app/auth"
	"scifi-search/app/database"
	"scifi-search/app/http/cookies"
	"scifi-search/app/http/middlewares"
	"strconv"

	"scifi-search/app/languages"
	"scifi-search/app/views"

	sqlc "scifi-search/app/database"
)

func RegisterFundingHandlers() {

	//http.HandleFunc("/funding", middlewares.AdminOnly(addFundingHandler))
	http.HandleFunc("/funding", middlewares.RequiresEmailVerified(middlewares.RequiresAuthorization(addFundingHandler, 1)))
	http.HandleFunc("/funding/update-items", middlewares.RequiresEmailVerified(middlewares.RequiresAuthorization(updateFundingItemsHandler, 1)))
}

// ------------------------------------------------------------------------------------------------

func addFundingHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		showFundingsManagementPage(w, r)
	case http.MethodPost:
		addFunding(w, r)
	case http.MethodDelete:
		//deleteFunding(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
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

func getFundingDocs(w http.ResponseWriter, r *http.Request, offset int) ([]map[string]any, int, error) {

	var fundings []map[string]any
	var fundingsDocs []database.Document

	user, err := auth.GetCurrentUser(w, r, queries)
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

	} else {

		// De ser loader, se listan solo sus documentos.
		fundingsDocs, err = queries.ListDocumentsByUser(r.Context(), sqlc.ListDocumentsByUserParams{
			UserID: sql.NullInt32{Int32: user.UserID, Valid: true},
			Limit:  10,
			Offset: int32((offset - 1) * 10),
		})
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

	totalFundings, err := queries.CountAllDocuments(r.Context())
	if err != nil {
		log.Fatal("Error al contar los documentos")
	}

	return fundings, int(totalFundings), nil
}

func updateFundingItemsHandler(w http.ResponseWriter, r *http.Request) {

	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1 // TODO: dejarlo así o levantamos error?
	}

	fundings, totalFundings, err := getFundingDocs(w, r, page)
	if err != nil {
		return
	}
	component := views.FundingList(fundings, page, totalFundings, languages.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func addFunding(w http.ResponseWriter, r *http.Request) {

	// Parsin del formulario.
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

	user, err := auth.GetCurrentUser(w, r, queries)
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
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	notifyFundingAddition(w, r, name)

	component := views.FundingAddedPage(languages.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func notifyFundingAddition(w http.ResponseWriter, r *http.Request, fundingName string) {

	// Acá sería mejor que el email solo se enviara a los usuarios
	// que están verificados.
	// TODO: agregar un campo a la base de datos que indique si
	// el usuario está verificado.
	// TODO: tampoco debería (creo) notificarse al usuario que añadió
	// el financiamiento.
	usersList, err := queries.ListUsers(r.Context())
	if err != nil {
		log.Fatal("Error en la notificación de nuevo financiamiento")
		return
	}

	for _, user := range usersList {

		userEmail := auth.GetUserEmail(user.AuthID)

		if userEmail == nil {
			log.Printf("Usuario %d no encontrado", user.UserID)
			continue
		}

		emailService.Send(*userEmail, "Nuevo financiamiento añadido", fundingName)
	}
}
