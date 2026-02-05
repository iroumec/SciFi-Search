package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"regexp"
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
	InvalidIDError         = errors.New("errors.invalid-ID")
	UserNotFoundError      = errors.New("errors.user-not-found")
	CountingDocumentsError = errors.New("errors.counting-documents")
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Endpoints.
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
		editFunding(w, r)
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

	component := views.ManageFundingPage(
		fundings,
		totalFundings,
		auth.GetCurrentAuthorizationLevel(w, r),
		documentTypes,
		documentAreas,
		documentCountriesBasedOn,
		documentGrantors,
		documentCurrencies,
		languages.GetTranslatorFromRequest(r),
	)
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func getFundingDocs(w http.ResponseWriter, r *http.Request, offset int) ([]map[string]any, int, error) {

	var fundings []map[string]any
	var fundingsDocs []database.Document
	var totalFundings int64

	user, err := getCurrentUser(w, r)
	if err != nil {
		cookies.AddFlashCookie(w, UnexpectedError.Error())
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, 0, err
	}

	if auth.GetAuthenticationLevel(user.AuthID) == auth.AdminRole.Level {

		// All documents are listed.
		fundingsDocs, err = queries.ListAllDocumentsWithLimitAndOffset(r.Context(),
			sqlc.ListAllDocumentsWithLimitAndOffsetParams{
				Limit:  10,
				Offset: int32((offset - 1) * 10),
			})

		totalFundings, err = queries.CountAllDocuments(r.Context())
		if err != nil {
			cookies.AddFlashCookie(w, CountingDocumentsError.Error())
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusInternalServerError)
			return nil, 0, err
		}

	} else {

		// Only the user's documents are listed.
		fundingsDocs, err = queries.ListDocumentsByUserWithLimitAndOffset(r.Context(),
			sqlc.ListDocumentsByUserWithLimitAndOffsetParams{
				UserID: sql.NullInt32{Int32: user.UserID, Valid: true},
				Limit:  10,
				Offset: int32((offset - 1) * 10),
			})

		totalFundings, err = queries.CountDocumentsByUser(r.Context(), sql.NullInt32{Int32: user.UserID, Valid: true})
		if err != nil {
			cookies.AddFlashCookie(w, CountingDocumentsError.Error())
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusInternalServerError)
			return nil, 0, err
		}
	}
	if err != nil {
		return nil, 0, err
	}

	for _, doc := range fundingsDocs {
		funding := map[string]any{
			"id":             doc.ID,
			"Name":           doc.Name,
			"Type":           doc.Type,
			"Main area":      doc.MainArea,
			"Secondary area": doc.SecondaryArea.String,
			"Link":           doc.Link.String,
			"Description":    doc.Description.String,
			"Based on":       doc.BasedOn.String,
			"Grantor":        doc.Grantor.String,
			"Currency":       doc.Currency,
			"Amount":         doc.Amount,
			"Deadline":       doc.Deadline,
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

	component := views.FundingList(
		fundings,
		page,
		totalFundings,
		languages.GetTranslatorFromRequest(r),
	)
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func getPage(r *http.Request) int {
	pageStr := r.URL.Query().Get("page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	return page
}

// ------------------------------------------------------------------------------------------------

func addFunding(w http.ResponseWriter, r *http.Request) {

	// Form parsing.
	if err := r.ParseForm(); err != nil {
		http.Error(w, FormParsingError.Error(), http.StatusBadRequest)
		return
	}

	// Form data obtention.
	name := r.Form.Get("name")
	fundingType := r.Form.Get("type")
	mainArea := r.Form.Get("main-area")
	secondaryArea := r.Form.Get("secondary-area")
	link := r.Form.Get("link")
	description := r.Form.Get("description")
	basedOn := r.Form.Get("based-on")
	grantor := r.Form.Get("grantor")
	currency := r.Form.Get("currency")
	amount := r.Form.Get("amount")
	deadline := r.Form.Get("deadline")

	translator := languages.GetTranslatorFromRequest(r)

	if err := validateFundingData(currency, amount); err != nil {
		notifications.ShowFlash(w, r, translator(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := getCurrentUser(w, r)
	if err != nil {
		cookies.AddFlashCookie(w, UnknownError.Error())
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	document, err := queries.AddDocument(r.Context(), sqlc.AddDocumentParams{
		Name:          name,
		UserID:        sql.NullInt32{Int32: user.UserID, Valid: true},
		Type:          fundingType,
		MainArea:      mainArea,
		SecondaryArea: sql.NullString{String: secondaryArea, Valid: secondaryArea != ""},
		Link:          sql.NullString{String: link, Valid: link != ""},
		Description:   sql.NullString{String: description, Valid: description != ""},
		BasedOn:       sql.NullString{String: basedOn, Valid: basedOn != ""},
		Grantor:       sql.NullString{String: grantor, Valid: grantor != ""},
		Currency:      currency,
		Amount:        amount,
		Deadline:      deadline,
	})

	_, err = client.Index(indexName).AddDocuments(map[string]any{
		"id":             document.ID,
		"User":           document.UserID,
		"Name":           document.Name,
		"Type":           document.Type,
		"Main area":      document.MainArea,
		"Secondary area": document.SecondaryArea.String,
		"Link":           document.Link.String,
		"Description":    document.Description.String,
		"Based on":       document.BasedOn.String,
		"Grantor":        document.Grantor.String,
		"Currency":       document.Currency,
		"Amount":         document.Amount,
		"Deadline":       document.Deadline,
	}, &primaryKey)
	if err != nil {
		log.Fatal(err)
	}

	notifyFundingAddition(r.Context(), document, translator)

	notifications.ShowFlash(w, r, translator("messages.funding-added"))

	updateFundingList(w, r)
}

// ------------------------------------------------------------------------------------------------

func notifyFundingAddition(
	ctx context.Context, document sqlc.Document, translator languages.Translator,
) {

	users, err := queries.UsersInterestedInFunding(ctx, document.ID)
	if err != nil {
		return
	}

	emailSubject := translator("messages.new-funding-email-subject")

	emailBody := strings.Join([]string{
		translator("messages.new-funding-email-body-part-1") + ": \n\n",
		translator("name") + ": " + document.Name + "\n\n",
		translator("description") + ": " + document.Description.String + "\n\n",
		translator("messages.new-funding-email-body-part-2") + "\n\n",
		translator("messages.new-funding-email-body-part-3"),
	}, "")

	for _, user := range users {
		notifyFundingAdditionViaEmail(user, emailSubject, emailBody)
	}
}

// ------------------------------------------------------------------------------------------------

func notifyFundingAdditionViaEmail(
	user database.UsersInterestedInFundingRow, emailSubject, emailBody string,
) {

	emailVerified, err := auth.IsEmailVerified(user.AuthID)
	if err != nil || !*emailVerified {
		return
	}

	userEmail := auth.GetUserEmail(user.AuthID)
	if userEmail == nil {
		return
	}

	emailService.Send(
		*userEmail,
		emailSubject,
		emailBody,
	)
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

	translator := languages.GetTranslatorFromRequest(r)

	notifications.ShowFlash(w, r, translator("messages.funding-deleted"))

	updateFundingList(w, r)
}

// ------------------------------------------------------------------------------------------------

func openAddNewFundingModal(w http.ResponseWriter, r *http.Request) {
	component := views.FundingModal(
		"new",
		nil,
		documentTypes,
		documentAreas,
		documentCountriesBasedOn,
		documentGrantors,
		documentCurrencies,
		languages.GetTranslatorFromRequest(r),
	)
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
		"id":             strconv.Itoa(int(sqlcDocument.ID)),
		"Name":           sqlcDocument.Name,
		"Type":           sqlcDocument.Type,
		"Main area":      sqlcDocument.MainArea,
		"Secondary area": sqlcDocument.SecondaryArea.String,
		"Link":           sqlcDocument.Link.String,
		"Description":    sqlcDocument.Description.String,
		"Based on":       sqlcDocument.BasedOn.String,
		"Grantor":        sqlcDocument.Grantor.String,
		"Currency":       sqlcDocument.Currency,
		"Amount":         sqlcDocument.Amount,
		"Deadline":       sqlcDocument.Deadline,
	}

	component := views.FundingModal(
		"edit",
		document,
		documentTypes,
		documentAreas,
		documentCountriesBasedOn,
		documentGrantors,
		documentCurrencies,
		languages.GetTranslatorFromRequest(r),
	)
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func editFunding(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/funding/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, InvalidIDError.Error(), http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, FormParsingError.Error(), http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	fundingType := r.Form.Get("type")
	mainArea := r.Form.Get("main-area")
	secondaryArea := r.Form.Get("secondary-area")
	link := r.Form.Get("link")
	description := r.Form.Get("description")
	basedOn := r.Form.Get("based-on")
	grantor := r.Form.Get("grantor")
	currency := r.Form.Get("currency")
	amount := r.Form.Get("amount")
	deadline := r.Form.Get("deadline")

	translator := languages.GetTranslatorFromRequest(r)

	if err := validateFundingData(currency, amount); err != nil {
		notifications.ShowFlash(w, r, translator(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := getCurrentUser(w, r)
	if err != nil {
		cookies.AddFlashCookie(w, UnknownError.Error())
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = queries.UpdateDocument(r.Context(), sqlc.UpdateDocumentParams{
		ID:            int32(id),
		Name:          name,
		UserID:        sql.NullInt32{Int32: user.UserID, Valid: true},
		Type:          fundingType,
		MainArea:      mainArea,
		SecondaryArea: sql.NullString{String: secondaryArea, Valid: secondaryArea != ""},
		Link:          sql.NullString{String: link, Valid: link != ""},
		Description:   sql.NullString{String: description, Valid: description != ""},
		BasedOn:       sql.NullString{String: basedOn, Valid: basedOn != ""},
		Grantor:       sql.NullString{String: grantor, Valid: grantor != ""},
		Currency:      currency,
		Amount:        amount,
		Deadline:      deadline,
	})
	if err != nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	_, err = client.Index(indexName).UpdateDocuments(map[string]any{
		"id":             id,
		"User":           user.UserID,
		"Name":           name,
		"Type":           fundingType,
		"Main area":      mainArea,
		"Secondary area": secondaryArea,
		"Link":           link,
		"Description":    description,
		"Based on":       basedOn,
		"Grantor":        grantor,
		"Currency":       currency,
		"Amount":         amount,
		"Deadline":       deadline,
	}, &primaryKey)
	if err != nil {
		log.Fatal(err)
	}

	notifications.ShowFlash(w, r, translator("messages.funding-edited"))

	updateFundingList(w, r)
}

// ------------------------------------------------------------------------------------------------

func validateFundingData(currency, amount string) error {
	// Currency format validation.
	currencyPattern := regexp.MustCompile(`^[A-Z]{3}$`)
	if !currencyPattern.MatchString(currency) {
		return errors.New("errors.invalid-currency-format")
	}

	// Amount format validation.
	amountPattern := regexp.MustCompile(`^[0-9]+([.,][0-9]+)?( - [0-9]+([.,][0-9]+)?)?$`)
	if !amountPattern.MatchString(amount) || !validateAmountRange(amount) {
		return errors.New("errors.invalid-amount-format")
	}

	return nil
}

// ------------------------------------------------------------------------------------------------

func validateAmountRange(amount string) bool {
	if strings.Contains(amount, " - ") {
		normalized := strings.ReplaceAll(amount, ",", ".")

		parts := strings.Split(normalized, " - ")
		if len(parts) != 2 {
			// Should never happen due to regex validation,
			// but added as an extra precaution.
			return false
		}
		minStr := strings.TrimSpace(parts[0])
		maxStr := strings.TrimSpace(parts[1])

		min, errMin := strconv.ParseFloat(minStr, 64)
		max, errMax := strconv.ParseFloat(maxStr, 64)
		if errMin != nil || errMax != nil {
			return false
		}
		return min < max
	}
	return true
}

// ------------------------------------------------------------------------------------------------
