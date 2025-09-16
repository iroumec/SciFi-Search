package providers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	sqlc "uki/app/database/sqlc"
)

const baseURL = "https://api.themoviedb.org/3"

type TMDbProvider struct {
	APIKey string
	Client *http.Client
}

func NewTMDbProvider(apiKey string) *TMDbProvider {
	return &TMDbProvider{
		APIKey: apiKey,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// ---------------------------------------------------
// MODELOS INTERNOS (estructuras que mapean la API)
// ---------------------------------------------------
type tmdbSearchResponse struct {
	Results []struct {
		ID         int    `json:"id"`
		Title      string `json:"title"`
		Name       string `json:"name"` // Para series
		Overview   string `json:"overview"`
		PosterPath string `json:"poster_path"`
		MediaType  string `json:"media_type"`
	} `json:"results"`
}

type tmdbMovie struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Overview   string `json:"overview"`
	PosterPath string `json:"poster_path"`
}

// ---------------------------------------------------
// IMPLEMENTACIÓN DE PROVIDER
// ---------------------------------------------------

func (p *TMDbProvider) Search(query string) ([]sqlc.Work, error) {
	endpoint := fmt.Sprintf("%s/search/multi?api_key=%s&query=%s", baseURL, p.APIKey, url.QueryEscape(query))

	resp, err := p.Client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data tmdbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var results []sqlc.Work
	for _, r := range data.Results {
		title := r.Title
		if title == "" {
			title = r.Name
		}

		typeID, ok := MapTMDbMediaTypeToID(r.MediaType)
		if !ok {
			continue // No se guardan los tipos que no interesan.
		}

		results = append(results, sqlc.Work{
			// ID:            fmt.Sprintf("%d", r.ID),  External ID
			Title:         title,
			ContentTypeID: int32(typeID),
			// Al ser la descripción un nullString, no puede asignarse directamente.
			// Se debe especificar un campo especificando si es válido o no.
			Description: convertToNullString(r.Overview),
			// ImageURL: fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", r.PosterPath),
			Source: "tmdb",
		})
	}

	return results, nil
}

func (p *TMDbProvider) GetByID(id string) (*sqlc.Work, error) {
	endpoint := fmt.Sprintf("%s/movie/%s?api_key=%s", baseURL, id, p.APIKey)

	resp, err := p.Client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var m tmdbMovie
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	content := &sqlc.Work{
		//ID:          fmt.Sprintf("%d", m.ID),
		Title:         m.Title,
		ContentTypeID: 0,
		Description:   convertToNullString(m.Overview),
		//ImageURL: fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", m.PosterPath),
		Source: "tmdb",
	}

	return content, nil
}

func convertToNullString(text string) sql.NullString {

	return sql.NullString{
		String: text,
		Valid:  text != "",
	}
}

func MapTMDbMediaTypeToID(mediaType string) (int, bool) {
	switch mediaType {
	case "movie":
		return 1, true // ID = 1
	case "tv":
		return 2, true // ID = 2
	default:
		return 0, false // Se ignoran otro tipos de resultados, como usuarios.
	}
}
