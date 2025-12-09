package handlers

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
	"scifi-search/app/utils"
	"scifi-search/app/views"
	"sort"
	"strings"

	sqlc "scifi-search/app/database"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

// ------------------------------------------------------------------------------------------------

var (
	client        meilisearch.ServiceManager
	documentTypes = list.New()
)

// ------------------------------------------------------------------------------------------------

const (
	indexName = "funding"
	dataPath  = "./resources/planillas/fundingRecords.json"
)

// ------------------------------------------------------------------------------------------------

type SearchResponse struct {
	Hits []any `json:"hits"`
}

// ------------------------------------------------------------------------------------------------

func registerSearchHandlers() {

	host := utils.GetEnv("MEILI_HOST", "http://meilisearch:7700")
	apiKey := utils.GetEnv("MEILI_API_KEY", "meili")

	// Creación del cliente de Meilisearch.
	client = meilisearch.New(host, meilisearch.WithAPIKey(apiKey))

	// Se indexan los datos.
	indexarDatos()

	// Se registra el handler.
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/funding", addFundingHandler)
}

// ------------------------------------------------------------------------------------------------

func indexarDatos() {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	var docs []map[string]any
	if err := json.Unmarshal(data, &docs); err != nil {
		log.Fatal(err)
	}

	var indexDocs []map[string]any
	for _, doc := range docs {

		nombre, ok := doc["Nombre"].(string)

		// Si el documento tiene un nombre válido.
		if ok {
			descripcion, _ := doc["Descripción"].(string)
			granArea1, _ := doc["Gran area 1"].(string)
			granArea2, _ := doc["Gran area 2"].(string)
			tipo, _ := doc["Tipo"].(string)
			link, _ := doc["Link"].(string)

			// Añadido del documento a la base de datos.
			document, err := queries.AddDocument(context.Background(), sqlc.AddDocumentParams{
				Name:        nombre,
				Description: descripcion,
				FirstArea:   granArea1,
				SecondArea:  sql.NullString{String: granArea2, Valid: granArea2 != ""},
				Type:        tipo,
				Link:        sql.NullString{String: link, Valid: link != ""},
			})
			if err != nil {
				log.Fatal(err)
			}

			// Indexado del documento.
			filtered := map[string]any{
				"id":          document.ID,
				"Nombre":      nombre,
				"Descripcion": descripcion,
				"Gran area 1": granArea1,
				"Gran area 2": granArea2,
				"Tipo":        tipo,
				"Link":        link,
			}

			indexDocs = append(indexDocs, filtered)

			utils.AddIfNotExists(documentTypes, tipo)
		}
	}

	index := client.Index(indexName)

	configureSearchSettings(index)
	configureFilterableAttributes(index)
	configureSortableAttributes(index)

	_, err = index.AddDocuments(indexDocs, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Datos indexados correctamente.")
}

// ------------------------------------------------------------------------------------------------

func configureSearchSettings(index meilisearch.IndexManager) {

	// Configuración de atributos en los que se busca.
	_, err := index.UpdateSearchableAttributes(&[]string{
		"Nombre",
		"Descripcion",
		"Gran area 1",
		"Gran area 2",
		"Tipo",
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

func handleSearch(w http.ResponseWriter, r *http.Request) {

	// Obtención de la query
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "missing query parameter 'query'", http.StatusBadRequest)
		return
	}

	// Obtención de parámetros de filtro y ordenamiento.
	// Se eliminan las comillas consecutivas al inicio y al final del parámetro.
	filterTipo := strings.Trim(r.URL.Query().Get("tipo"), `"`)
	filterArea := strings.Trim(r.URL.Query().Get("area"), `"`)
	sortBy := strings.Trim(r.URL.Query().Get("sort"), `"`)

	// Construcción de la búsqueda.
	searchRequest := &meilisearch.SearchRequest{
		Limit:            20,
		ShowRankingScore: true, // Se muestra el score de relevancia.
	}

	// Aplicación de filtros.
	var filters []string
	if filterTipo != "" {
		filters = append(filters, fmt.Sprintf("Tipo = '%s'", filterTipo))
	}
	if filterArea != "" {
		filters = append(filters, fmt.Sprintf("\"Gran area 1\" = '%s' OR \"Gran area 2\" = '%s'", filterArea, filterArea))
	}

	if len(filters) > 0 {
		searchRequest.Filter = filters
	}

	// Aplicación de ordenamiento.
	/*if sortBy != "" {
		searchRequest.Sort = []string{sortBy}
	}*/
	// El anterior es un ordenamiento previo a la relevancia.

	res, err := client.Index(indexName).Search(query, searchRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// No puedo utilizar res.Hist directamnete porque es un slice reservado
	// de Meilisearch. Debo almacenarlo en una variable local.
	hits := make([]any, len(res.Hits))
	for i, h := range res.Hits {
		hits[i] = h
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

	// Ordenar los resultados después de obtenerlos (para mantener relevancia base).
	// TODO: llevarlo a una opción extra.
	if sortBy != "" {
		sortResults(hitsMaps, sortBy)
	}

	// De estar autenticado el usuario, se guarda la búsqueda
	// en su historial.
	user := getCurrentUser(w, r)
	if user != nil {
		params := sqlc.CreateHistoricSearchParams{UserID: user.UserID, SearchString: query}
		queries.CreateHistoricSearch(r.Context(), params)
	}

	// Pasar maps al templ.
	component := views.SearchResultsPage(query, hitsMaps, isUserAuthenticated(r), documentTypes)
	component.Render(r.Context(), w)

}

// ------------------------------------------------------------------------------------------------

func addFundingHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		showAddFundingPage(w, r)
	case http.MethodPost:
		addFunding(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func showAddFundingPage(w http.ResponseWriter, r *http.Request) {

	component := views.AddFundingPage()
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
	description := r.Form.Get("description")
	firstArea := r.Form.Get("firsrt-area")
	secondArea := r.Form.Get("second-area")
	fundingType := r.Form.Get("type")
	link := r.Form.Get("link")

	document, err := queries.AddDocument(r.Context(), sqlc.AddDocumentParams{
		Name:        name,
		Description: description,
		FirstArea:   firstArea,
		SecondArea:  sql.NullString{String: secondArea, Valid: secondArea != ""},
		Type:        fundingType,
		Link:        sql.NullString{String: link, Valid: link != ""},
	})

	_, err = client.Index(indexName).AddDocuments(map[string]any{
		"id":          document.ID,
		"Nombre":      document.Name,
		"Descripcion": document.Description,
		"Gran area 1": document.FirstArea,
		"Gran area 2": document.SecondArea.String,
		"Tipo":        document.Type,
		"Link":        document.Link.String,
	}, nil)
	if err != nil {
		log.Fatal(err)
	}

	component := views.FundingAddedPage(utils.GetTranslatorFromRequest(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------

func configureFilterableAttributes(index meilisearch.IndexManager) {

	_, err := index.UpdateFilterableAttributes(&[]any{
		"Tipo",
		"Gran area 1",
		"Gran area 2",
	})
	if err != nil {
		log.Println("Error configurando filtros:", err)
	}
}

// ------------------------------------------------------------------------------------------------

func configureSortableAttributes(index meilisearch.IndexManager) {

	configureRankingRules(index)

	_, err := index.UpdateSortableAttributes(&[]string{
		"Nombre",
		"Tipo",
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
func sortResults(hits []map[string]any, sortBy string) {
	// Parseo del sortBy (ej: "Nombre:asc" o "Tipo:desc").
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
