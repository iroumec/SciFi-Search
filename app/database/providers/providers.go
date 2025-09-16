package providers

import (
	sqlc "uki/app/database/sqlc"
)

// Interfaz para cualquier proveedor externo.
type Provider interface {
	Search(query string) ([]sqlc.Work, error)
	GetByID(id string) (*sqlc.Work, error)
}

// Repositorio unififcado.
type Repository struct {
	Providers []Provider
	Cache     sqlc.Queries // tu DB interna
}

func (r *Repository) Search(query string) ([]sqlc.Work, error) {
	var results []sqlc.Work
	for _, p := range r.Providers {
		data, err := p.Search(query)
		if err == nil {
			results = append(results, data...)
			// Guardado en nuestra base de datos como caché.
			//r.Cache.Save(data)
		}
	}
	return results, nil
}
