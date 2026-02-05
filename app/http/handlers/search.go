package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"container/list"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"scifi-search/app/auth"
	"scifi-search/app/languages"
	"scifi-search/app/utils"
	"scifi-search/app/utils/structures"
	"scifi-search/app/views"
	"sort"
	"strconv"
	"strings"

	sqlc "scifi-search/app/database"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

var (
	client                   meilisearch.ServiceManager
	documentTypes            = list.New()
	documentAreas            = list.New()
	documentCountriesBasedOn = list.New()
	documentGrantors         = list.New()
	documentCurrencies       = list.New()
	primaryKey               = "id"
)

// Errors.
var (
	SearchableAttributesConfigurationError = errors.New(
		"errors.searchable-attributes-configuration",
	)
	TypoToleranceConfigurationError = errors.New("errors.typo-tolerance-configuration")
	SynonymsConfigurationError      = errors.New("errors.synonyms-configuration")
	FiltersConfigurationError       = errors.New("errors.filters-configuration")
	OrderingConfigurationError      = errors.New("errors.ordering-configuration")
	RankingConfigurationError       = errors.New("errors.ranking-configuration")
	MissingQueryParameterError      = errors.New("errors.missing-query-parameter")
)

// ------------------------------------------------------------------------------------------------
// Constants
// ------------------------------------------------------------------------------------------------

const (
	indexName      = "funding"
	dataPath       = "./resources/planillas/fundingRecords.json"
	resultsPerPage = 10
)

// ------------------------------------------------------------------------------------------------
// Structures
// ------------------------------------------------------------------------------------------------

type SearchResponse struct {
	Hits []any `json:"hits"`
}

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

// Endpoints.
func RegisterSearchHandlers() {

	host := utils.GetEnv("MEILI_HOST", "http://meilisearch:7700")
	apiKey := utils.GetEnv("MEILI_API_KEY", "meili")

	// Client creation.
	client = meilisearch.New(host, meilisearch.WithAPIKey(apiKey))

	// Data indexation.
	indexData()

	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/search/update-filter", filtersHandler)
	http.HandleFunc("/search/update-results", filtersHandler)
}

// ------------------------------------------------------------------------------------------------
// Definición de handlers
// ------------------------------------------------------------------------------------------------

func searchHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		showSearchResults(w, r)
	default:
		http.Error(w, MethodNotAllowedError.Error(), http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func filtersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		updateResults(w, r)
	default:
		http.Error(w, MethodNotAllowedError.Error(), http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func indexData() {

	// In case data has already been indexed...
	if indexContainsDocuments(indexName) {
		return
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	var documents []map[string]any
	if err := json.Unmarshal(data, &documents); err != nil {
		log.Fatal(err)
	}

	indexDocs := indexDocuments(documents)

	index := client.Index(indexName)

	configureSearchSettings(index)
	configureFilterableAttributes(index)
	configureSortableAttributes(index)

	_, err = index.AddDocuments(indexDocs, &primaryKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Successful indexation!")
}

// ------------------------------------------------------------------------------------------------

func indexDocuments(documents []map[string]any) []map[string]any {

	var indexDocs []map[string]any
	for _, doc := range documents {

		name, ok := doc["name"].(string)

		// If the document name is valid...
		if ok {
			documentType, _ := doc["type"].(string)
			mainArea, _ := doc["main_area"].(string)
			secondaryArea, _ := doc["secondary_area"].(string)
			link, _ := doc["link"].(string)
			description, _ := doc["description"].(string)
			basedOn, _ := doc["based_on"].(string)
			grantor, _ := doc["grantor"].(string)
			currency, _ := doc["currency"].(string)
			amount, _ := doc["amount"].(string)
			deadline, _ := doc["deadline"].(string)

			// The document is added to the database.
			document, err := queries.AddDocument(context.Background(), sqlc.AddDocumentParams{
				Name:          name,
				UserID:        sql.NullInt32{Valid: false}, // Undefined.
				Type:          documentType,
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
				log.Fatal(err)
			}

			// Document indexing.
			filtered := map[string]any{
				"id":             document.ID,
				"User":           document.UserID,
				"Name":           name,
				"Type":           documentType,
				"Main area":      mainArea,
				"Secondary area": secondaryArea,
				"Link":           link,
				"Description":    description,
				"Based on":       basedOn,
				"Grantor":        grantor,
				"Currency":       currency,
				"Amount":         amount,
				"Deadline":       deadline,
			}

			indexDocs = append(indexDocs, filtered)

			// These lists are solely used for autocomplete fields in the manage funding page.
			structures.AddIfNotExists(documentTypes, documentType)
			structures.AddIfNotExists(documentAreas, mainArea)
			structures.AddIfNotExists(documentAreas, secondaryArea)
			structures.AddIfNotExists(documentCountriesBasedOn, basedOn)
			structures.AddIfNotExists(documentGrantors, grantor)
			structures.AddIfNotExists(documentCurrencies, currency)
		}
	}

	return indexDocs
}

// ------------------------------------------------------------------------------------------------

func indexContainsDocuments(indexName string) bool {

	index := client.Index(indexName)

	stats, err := index.GetStats()
	if err != nil {
		// In case the index doesn't exist...
		if strings.Contains(err.Error(), "index_not_found") {
			return false
		}
		log.Fatal(UnknownError)
	}

	return stats.NumberOfDocuments > 0
}

// ------------------------------------------------------------------------------------------------

func configureSearchSettings(index meilisearch.IndexManager) {

	// Searchable attributes.
	_, err := index.UpdateSearchableAttributes(&[]string{
		"Name",
		"Type",
		"Main area",
		"Secondary area",
		"Description",
		"Grantor",
	})
	if err != nil {
		log.Println(SearchableAttributesConfigurationError)
	}

	// Typo tolerance.
	_, err = index.UpdateTypoTolerance(&meilisearch.TypoTolerance{
		Enabled: true,
		MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
			OneTypo:  4, // 4+ letters words: 1 typo allowed.
			TwoTypos: 8, // 8+ letters words: 2 typos allowed.
		},
	})
	if err != nil {
		log.Println(TypoToleranceConfigurationError)
	}

	// Commun synonyms.
	_, err = index.UpdateSynonyms(&map[string][]string{
		"engineering": {"ingenieria", "ingeniería"},
		"ingenieria":  {"engineering", "ingeniería"},
		"science":     {"ciencia", "ciencias"},
		"ciencia":     {"science", "ciencias"},
		"tech":        {"technology", "tecnologia", "tecnología"},
		"tecnologia":  {"tech", "technology"},
	})
	if err != nil {
		log.Println(SynonymsConfigurationError)
	}
}

// ------------------------------------------------------------------------------------------------

func showSearchResults(w http.ResponseWriter, r *http.Request) {

	translator := languages.GetTranslatorFromRequest(r)

	// Query obtention.
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, translator(MissingQueryParameterError.Error()), http.StatusBadRequest)
		return
	}

	// Search construction.
	searchRequest := &meilisearch.SearchRequest{
		ShowRankingScore: true, // Relevance score.
		Limit:            1000,
	}

	res, err := client.Index(indexName).Search(query, searchRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hits := make([]any, 0, resultsPerPage)
	for i := 0; i < len(res.Hits) && i < resultsPerPage; i++ {
		hits = append(hits, res.Hits[i])
	}

	// []map[string]any convetion.
	data, err := json.Marshal(hits)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var hitsMaps []map[string]any
	if err := json.Unmarshal(data, &hitsMaps); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// These lists are solely used for filters in the search results page.
	types, areas, countriesBasedOn := getFilterValues(hitsMaps)

	var authorizationLevel = auth.NoRole.Level

	// If the user is authenticated, the search is saved in their history.
	if auth.IsUserAuthenticated(w, r) {

		user, err := getCurrentUser(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authorizationLevel = auth.GetAuthenticationLevel(user.AuthID)

		params := sqlc.CreateHistoricSearchParams{UserID: user.UserID, SearchString: query}
		queries.CreateHistoricSearch(r.Context(), params)
	}

	component := views.SearchResultsPage(
		query,
		hitsMaps,
		len(res.Hits),
		types,
		areas,
		countriesBasedOn,
		authorizationLevel,
		languages.GetTranslatorFromRequest(r),
	)
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func updateResults(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("query")
	filterBasedOn := r.URL.Query()["based-on"]
	filterType := r.URL.Query()["type"]
	filterArea := r.URL.Query()["area"]
	sortBy := r.URL.Query()["sortby"]
	pageStr := r.URL.Query().Get("page")

	var filters []string

	for _, t := range filterBasedOn {
		filters = append(filters, fmt.Sprintf("\"Based on\" = '%s'", t))
	}

	for _, t := range filterType {
		filters = append(filters, fmt.Sprintf("Type = '%s'", t))
	}

	for _, t := range filterArea {
		filters = append(filters,
			fmt.Sprintf("\"Main area\" = '%s' OR \"Secondary area\" = '%s'", t, t),
		)
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	fullSearchRequest := &meilisearch.SearchRequest{
		Limit: 1000,
	}

	searchRequest := &meilisearch.SearchRequest{
		ShowRankingScore: true,
		Limit:            resultsPerPage,
		Offset:           int64(page-1) * resultsPerPage,
	}

	allHits := getFilteredResponse(w, filters, sortBy, query, fullSearchRequest)
	res := getFilteredResponse(w, filters, sortBy, query, searchRequest)

	// No results.
	if res == nil || len(res.Hits) == 0 {
		component := views.SearchResults(
			[]map[string]any{},
			0,
			page,
			languages.GetTranslatorFromRequest(r),
		)
		component.Render(r.Context(), w)
		return
	}

	// Dynamic slice construction.
	hits := make([]any, 0, len(res.Hits))
	for _, h := range res.Hits {
		hits = append(hits, h)
	}

	// []map[string]any convetions.
	data, err := json.Marshal(hits)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var hitsMaps []map[string]any
	if err := json.Unmarshal(data, &hitsMaps); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ordenamiento.
	if len(sortBy) > 0 {
		sortResults(hitsMaps, sortBy)
	}

	total := 0
	if allHits != nil {
		total = len(allHits.Hits)
	}

	component := views.SearchResults(hitsMaps, total, page, languages.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func getFilteredResponse(
	w http.ResponseWriter,
	filters []string,
	sortBy []string,
	query string,
	searchRequest *meilisearch.SearchRequest,
) *meilisearch.SearchResponse {
	if len(filters) > 0 {
		searchRequest.Filter = filters
	}

	// Ordering application.
	if len(sortBy) > 0 {
		searchRequest.Sort = sortBy
	}
	// This ordering is previous to the relevance.

	res, err := client.Index(indexName).Search(query, searchRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return res
}

// ------------------------------------------------------------------------------------------------

func configureFilterableAttributes(index meilisearch.IndexManager) {

	_, err := index.UpdateFilterableAttributes(&[]any{
		"Type",
		"Main area",
		"Secondary area",
		"Based on",
	})
	if err != nil {
		log.Println(FiltersConfigurationError)
	}
}

// ------------------------------------------------------------------------------------------------

func configureSortableAttributes(index meilisearch.IndexManager) {

	configureRankingRules(index)

	_, err := index.UpdateSortableAttributes(&[]string{
		"Name",
		"Type",
		"Deadline",
	})
	if err != nil {
		log.Println(OrderingConfigurationError)
	}
}

// ------------------------------------------------------------------------------------------------

// Ranking rules configuration to maintain relevance.
func configureRankingRules(index meilisearch.IndexManager) {

	_, err := index.UpdateRankingRules(&[]string{
		"words",
		"typo",
		"proximity",
		"attribute",
		"sort", // Sort  applied affter relevance.
		"exactness",
	})
	if err != nil {
		log.Println(RankingConfigurationError)
	}
}

// ------------------------------------------------------------------------------------------------

func sortResults(hits []map[string]any, sortByArray []string) {
	sortBy := sortByArray[0] // It always has a unique element.

	// `sortBy` parsing (for example, "Name:asc" o "Type:desc").
	parts := splitSort(sortBy)
	if len(parts) != 2 {
		return
	}

	field := parts[0]
	order := parts[1]

	// Ordering.
	sort.Slice(hits, func(i, j int) bool {
		valI, okI := hits[i][field].(string)
		valJ, okJ := hits[j][field].(string)

		if !okI || !okJ {
			return false
		}

		if order == "asc" {
			return valI < valJ
		}
		return valI > valJ
	})
}

// ------------------------------------------------------------------------------------------------

// It divides the sort parameter.
func splitSort(sortBy string) []string {
	for i, char := range sortBy {
		if char == ':' {
			return []string{sortBy[:i], sortBy[i+1:]}
		}
	}
	return []string{sortBy}
}

// ------------------------------------------------------------------------------------------------

func getFilterValues(hitsMaps []map[string]any) (*list.List, *list.List, *list.List) {
	documentTypes := list.New()
	documentAreas := list.New()
	documentCountriesBasedOn := list.New()

	for _, doc := range hitsMaps {
		if documentType, ok := doc["Type"].(string); ok {
			structures.AddIfNotExists(documentTypes, documentType)
		}
		if mainArea, ok := doc["Main area"].(string); ok {
			structures.AddIfNotExists(documentAreas, mainArea)
		}
		if secondaryArea, ok := doc["Secondary area"].(string); ok {
			structures.AddIfNotExists(documentAreas, secondaryArea)
		}
		if basedOn, ok := doc["Based on"].(string); ok {
			structures.AddIfNotExists(documentCountriesBasedOn, basedOn)
		}
	}

	return documentTypes, documentAreas, documentCountriesBasedOn
}

// ------------------------------------------------------------------------------------------------
