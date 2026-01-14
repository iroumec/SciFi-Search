package handlers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"container/list"
	"context"
	"database/sql"
	"encoding/json"
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

// ------------------------------------------------------------------------------------------------
// Constantes
// ------------------------------------------------------------------------------------------------

const (
	indexName      = "funding"
	dataPath       = "./resources/planillas/fundingRecords.json"
	resultsPerPage = 10
)

// ------------------------------------------------------------------------------------------------
// Estructuras
// ------------------------------------------------------------------------------------------------

type SearchResponse struct {
	Hits []any `json:"hits"`
}

// ------------------------------------------------------------------------------------------------
// Registro de endpoints
// ------------------------------------------------------------------------------------------------

func RegisterSearchHandlers() {

	host := utils.GetEnv("MEILI_HOST", "http://meilisearch:7700")
	apiKey := utils.GetEnv("MEILI_API_KEY", "meili")

	// Creación del cliente de Meilisearch.
	client = meilisearch.New(host, meilisearch.WithAPIKey(apiKey))

	// Se indexan los datos.
	indexarDatos()

	// Se registra el handler.
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
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func filtersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		updateResults(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func indexarDatos() {

	// Si los datos ya fueron indexados, se retorna...
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

	fmt.Println("Datos indexados correctamente.")
}

// ------------------------------------------------------------------------------------------------

func indexDocuments(documents []map[string]any) []map[string]any {

	var indexDocs []map[string]any
	for _, doc := range documents {

		name, ok := doc["name"].(string)

		// Si el documento tiene un nombre válido.
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

			// Añadido del documento a la base de datos.
			document, err := queries.AddDocument(context.Background(), sqlc.AddDocumentParams{
				Name:          name,
				UserID:        sql.NullInt32{Valid: false}, // Indefinido en este caso.
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

			// Indexado del documento.
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
		// Si el índice no existe...
		if strings.Contains(err.Error(), "index_not_found") {
			return false
		}
		log.Fatal(err)
	}

	return stats.NumberOfDocuments > 0
}

// ------------------------------------------------------------------------------------------------

func configureSearchSettings(index meilisearch.IndexManager) {

	// Configuración de atributos en los que se busca.
	_, err := index.UpdateSearchableAttributes(&[]string{
		"Name",
		"Type",
		"Main area",
		"Secondary area",
		"Description",
		"Grantor",
	})
	if err != nil {
		log.Println("Error configurando atributos de búsqueda:", err)
	}

	// Configuración de typo tolerance (tolerancia a errores tipográficos).
	_, err = index.UpdateTypoTolerance(&meilisearch.TypoTolerance{
		Enabled: true,
		MinWordSizeForTypos: meilisearch.MinWordSizeForTypos{
			OneTypo:  4, // Palabras de 4+ letras: 1 error permitido.
			TwoTypos: 8, // Palabras de 8+ letras: 2 errores permitidos.
		},
	})
	if err != nil {
		log.Println("Error configurando tolerancia de typos:", err)
	}

	// Configuración de sinónimos comunes.
	_, err = index.UpdateSynonyms(&map[string][]string{
		"engineering": {"ingenieria", "ingeniería"},
		"ingenieria":  {"engineering", "ingeniería"},
		"science":     {"ciencia", "ciencias"},
		"ciencia":     {"science", "ciencias"},
		"tech":        {"technology", "tecnologia", "tecnología"},
		"tecnologia":  {"tech", "technology"},
	})
	if err != nil {
		log.Println("Error configurando sinónimos:", err)
	}
}

// ------------------------------------------------------------------------------------------------

func showSearchResults(w http.ResponseWriter, r *http.Request) {

	// Obtención de la query
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "missing query parameter 'query'", http.StatusBadRequest)
		return
	}

	// Construcción de la búsqueda.
	searchRequest := &meilisearch.SearchRequest{
		ShowRankingScore: true, // Se muestra el score de relevancia.
		Limit:            1000,
	}

	res, err := client.Index(indexName).Search(query, searchRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// No puedo utilizar res.Hist directamnete porque es un slice reservado
	// de Meilisearch. Debo almacenarlo en una variable local.
	hits := make([]any, 0, resultsPerPage)
	for i := 0; i < len(res.Hits) && i < resultsPerPage; i++ {
		hits = append(hits, res.Hits[i])
	}

	// Convertir a []map[string]any de forma segura
	data, err := json.Marshal(hits)
	if err != nil {
		log.Println("Error marshal hits:", err)
	}

	var hitsMaps []map[string]any
	if err := json.Unmarshal(data, &hitsMaps); err != nil {
		log.Println("Error unmarshal hits:", err)
	}

	// These lists are solely used for filters in the search results page.
	types, areas, countriesBasedOn := getFilterValues(hitsMaps)

	var authorizationLevel = auth.NoRole.Level

	// De estar autenticado el usuario, se guarda la búsqueda
	// en su historial.
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

	// Pasar maps al templ.
	component := views.SearchResultsPage(query, hitsMaps, len(res.Hits), types, areas, countriesBasedOn, authorizationLevel, languages.GetTranslatorFromRequest(r))
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

	log.Println(sortBy)

	var filters []string

	for _, t := range filterBasedOn {
		filters = append(filters, fmt.Sprintf("Based on = '%s'", t))
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

	// Sin resultados reales.
	if res == nil || len(res.Hits) == 0 {
		component := views.SearchResults([]map[string]any{}, 0, page, languages.GetTranslatorFromRequest(r))
		component.Render(r.Context(), w)
		return
	}

	// Construcción dinámica del slice.
	hits := make([]any, 0, len(res.Hits))
	for _, h := range res.Hits {
		hits = append(hits, h)
	}

	// Conversión a []map[string]any.
	data, err := json.Marshal(hits)
	if err != nil {
		log.Println("Error marshal hits:", err)
	}

	var hitsMaps []map[string]any
	if err := json.Unmarshal(data, &hitsMaps); err != nil {
		log.Println("Error unmarshal hits:", err)
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

func getFilteredResponse(w http.ResponseWriter, filters []string, sortBy []string, query string, searchRequest *meilisearch.SearchRequest) *meilisearch.SearchResponse {
	if len(filters) > 0 {
		searchRequest.Filter = filters
	}

	// Aplicación de ordenamiento.
	if len(sortBy) > 0 {
		searchRequest.Sort = sortBy
	}
	// El anterior es un ordenamiento previo a la relevancia.

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
		log.Println("Error configurando filtros:", err)
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
		log.Println("Error configurando ordenamiento:", err)
	}
}

// ------------------------------------------------------------------------------------------------

// Configuración de las ranking rules para mantener relevancia.
func configureRankingRules(index meilisearch.IndexManager) {

	_, err := index.UpdateRankingRules(&[]string{
		"words",
		"typo",
		"proximity",
		"attribute",
		"sort", // El sort se aplica después de la relevancia.
		"exactness",
	})
	if err != nil {
		log.Println("Error configurando ranking:", err)
	}
}

// ------------------------------------------------------------------------------------------------

// Función auxiliar para ordenar resultados
func sortResults(hits []map[string]any, sortByArray []string) {
	sortBy := sortByArray[0] //el array siempre tiene un solo elemento

	// Parseo del sortBy (ej: "Name:asc" o "Type:desc").
	parts := splitSort(sortBy)
	if len(parts) != 2 {
		return
	}

	field := parts[0]
	order := parts[1]

	// Ordenamiento.
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

// Función auxiliar para dividir el parámetro sort.
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
