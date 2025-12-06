package handlers

// ------------------------------------------------------------------------------------------------

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"scifi-search/app/utils"
	"scifi-search/app/views"

	sqlc "scifi-search/app/database"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

// ------------------------------------------------------------------------------------------------

var client meilisearch.ServiceManager

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
		}
	}

	index := client.Index(indexName)
	_, err = index.AddDocuments(indexDocs, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Datos indexados correctamente.")
}

// ------------------------------------------------------------------------------------------------

func handleSearch(w http.ResponseWriter, r *http.Request) {

	// Obtención de la query
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "missing query parameter 'query'", http.StatusBadRequest)
		return
	}

	res, err := client.Index(indexName).Search(query, &meilisearch.SearchRequest{})
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

	// De estar autenticado el usuario, se guarda la búsqueda
	// en su historial.
	user := getCurrentUser(w, r)
	if user != nil {
		params := sqlc.CreateHistoricSearchParams{UserID: user.UserID, SearchString: query}
		queries.CreateHistoricSearch(r.Context(), params)
	}

	// Pasar maps al templ.
	component := views.SearchResultsPage(query, hitsMaps, isUserAuthenticated(r))
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

	component := views.AddFundingPage(isUserAuthenticated(r))
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

	component := views.FundingAddedPage(isUserAuthenticated(r))
	component.Render(r.Context(), w)
}

// ------------------------------------------------------------------------------------------------
