package handlers

import (
	"encoding/json"
	"net/http"
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
