package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	meilisearch "github.com/meilisearch/meilisearch-go"
)

//var client *meilisearch.Client

func init() {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://meilisearch:7700",
		APIKey: getEnv("MEILI_KEY"),
	})

	index := client.Index("productos")
	res, err := index.Search("celular", &meilisearch.SearchRequest{
		Limit: 10,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Hits)
}

func buscarHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query param 'q' is required", http.StatusBadRequest)
		return
	}

	index := client.Index("productos")
	res, err := index.Search(query, &meilisearch.SearchRequest{
		Limit: 10,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res.Hits)
}

/*
func main() {
	http.HandleFunc("/api/buscar", buscarHandler)
	log.Println("Servidor Go escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
*/
