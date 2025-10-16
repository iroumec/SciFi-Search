package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"tpe/web/app/utils"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

// ------------------------------------------------------------------------------------------------

var client meilisearch.ServiceManager

// ------------------------------------------------------------------------------------------------

type SearchResponse struct {
	Hits []any `json:"hits"`
}

// ------------------------------------------------------------------------------------------------

func registerSearchHandlers() {

	initMeilisearch()

	http.HandleFunc("/search", handleSearch)
}

// ------------------------------------------------------------------------------------------------

func initMeilisearch() {

	host := utils.GetEnv("MEILI_HOST", "http://meilisearch:7700")
	apiKey := utils.GetEnv("MEILI_API_KEY", "meili")

	client = meilisearch.New(host, meilisearch.WithAPIKey(apiKey))

	indexarDatos()
}

func indexarDatos() {
	data, err := os.ReadFile("resources/planillas/fundingRecords.json")
	if err != nil {
		log.Fatal(err)
	}

	var docs []map[string]interface{}
	if err := json.Unmarshal(data, &docs); err != nil {
		log.Fatal(err)
	}

	_, err = client.Index("funding").AddDocuments(docs, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Datos indexados correctamente.")
}

// ------------------------------------------------------------------------------------------------

func handleSearch(w http.ResponseWriter, r *http.Request) {
	// Obtención de la query
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	res, err := client.Index("usuarios").Search(query, &meilisearch.SearchRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hits := make([]interface{}, len(res.Hits))
	for i, h := range res.Hits {
		hits[i] = h
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SearchResponse{Hits: hits})
}

// ------------------------------------------------------------------------------------------------
