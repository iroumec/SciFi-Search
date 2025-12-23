package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"scifi-search/app/http/cookies"
	"scifi-search/app/infra/email"

	//"scifi-search/app/http/middlewares"
	"scifi-search/app/languages"
	"scifi-search/app/views"

	sqlc "scifi-search/app/database"
)

func registerFundingHandlers() {

	//http.HandleFunc("/funding", middlewares.AdminOnly(addFundingHandler))
	http.HandleFunc("/funding", addFundingHandler)

}

// ------------------------------------------------------------------------------------------------

func addFundingHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		//showAddFundingPage(w, r)
		views.ManageFundingPage(languages.GetTranslatorFromRequest(r)).Render(r.Context(), w)
	case http.MethodPost:
		addFunding(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func showAddFundingPage(w http.ResponseWriter, r *http.Request) {

	if !isEmailVerified(w, r) {

		cookies.AddFlashCookie(w, languages.GetTranslatorFromRequest(r)("Debe verificar su email antes de acceder a esta funcionalidad."))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	component := views.AddFundingPage(languages.GetTranslatorFromRequest(r))
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
	deadline := r.Form.Get("deadline")

	document, err := queries.AddDocument(r.Context(), sqlc.AddDocumentParams{
		Name:        name,
		UserID:      getCurrentUser(w, r).UserID,
		Type:        fundingType,
		FirstArea:   firstArea,
		SecondArea:  sql.NullString{String: secondArea, Valid: secondArea != ""},
		Link:        sql.NullString{String: link, Valid: link != ""},
		Description: sql.NullString{String: description, Valid: description != ""},
		BasedOn:     sql.NullString{String: basedOn, Valid: basedOn != ""},
		Grantor:     sql.NullString{String: grantor, Valid: grantor != ""},
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
	users, err := queries.ListUsers(r.Context())
	if err != nil {
		log.Fatal("Error en la notificación de nuevo financiamiento")
		return
	}

	for _, user := range users {

		userEmail := getUserEmail(user.AuthID)

		if userEmail == nil {
			log.Printf("Usuario %d no encontrado", user.UserID)
			continue
		}

		email.Send(*userEmail, "Nuevo financiamiento añadido", fundingName)
	}
}
