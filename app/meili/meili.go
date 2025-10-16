package meili

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"tpe/web/app/utils"

	sqlc "tpe/web/app/database"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

// ------------------------------------------------------------------------------------------------

var client meilisearch.ServiceManager

// ------------------------------------------------------------------------------------------------

var queries *sqlc.Queries

// ------------------------------------------------------------------------------------------------

type SearchResponse struct {
	Hits []any `json:"hits"`
}

// ------------------------------------------------------------------------------------------------

func Init(q *sqlc.Queries) {

	queries = q

	host := utils.GetEnv("MEILI_HOST", "http://meilisearch:7700")
	apiKey := utils.GetEnv("MEILI_API_KEY", "meili")

	client = meilisearch.New(host, meilisearch.WithAPIKey(apiKey))

	indexarDatos()

	http.HandleFunc("/search", handleSearch)
}

// ------------------------------------------------------------------------------------------------

func indexarDatos() {
	data, err := os.ReadFile("resources/planillas/fundingRecords.json")
	if err != nil {
		log.Fatal(err)
	}

	var docs []map[string]any
	if err := json.Unmarshal(data, &docs); err != nil {
		log.Fatal(err)
	}

	var indexDocs []map[string]any

	for i, doc := range docs {
		// Crear id único
		if doc["ofibusubID"] == nil {
			doc["id"] = fmt.Sprintf("doc-%d", i)
		} else {
			doc["id"] = doc["ofibusubID"]
		}

		// Solo los campos relevantes
		filtered := map[string]any{
			"id":          doc["id"],
			"Nombre":      doc["Nombre"],
			"Descripcion": doc["Descripcion"],
			"Gran area 1": doc["Gran area 1"],
			"Gran area 2": doc["Gran area 2"],
			"Tipo":        doc["Tipo"],
		}

		indexDocs = append(indexDocs, filtered)
	}

	index := client.Index("funding")

	_, err = index.AddDocuments(indexDocs, nil) /*, &meilisearch.Add{
		PrimaryKey: "id",
	})*/
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

	res, err := client.Index("funding").Search(query, &meilisearch.SearchRequest{})
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SearchResponse{Hits: hits})
	params := sqlc.CreateHistoricSearchParams{UserID: 1, SearchString: query}
	queries.CreateHistoricSearch(r.Context(), params)
}

// ------------------------------------------------------------------------------------------------
