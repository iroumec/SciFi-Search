package handlers

import (
	"errors"
	"scifi-search/app/avatars"
	sqlc "scifi-search/app/database"
	"scifi-search/app/email"
)

// ------------------------------------------------------------------------------------------------
// Constantes
// ------------------------------------------------------------------------------------------------

// Ruta a partir de la cual se servirán los archivos estáticos.
const (
	debug = false
)

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

// Dependencias.
var (
	queries       *sqlc.Queries
	avatarService *avatars.Service
	emailService  *email.Service
)

// ------------------------------------------------------------------------------------------------

// Errores.
var (
	MethodNotAllowedError = errors.New("Method not allowed")
)

// ------------------------------------------------------------------------------------------------
// Setters
// ------------------------------------------------------------------------------------------------

func SetQueries(q *sqlc.Queries) {
	queries = q
}

// ------------------------------------------------------------------------------------------------

func SetAvatarService(service *avatars.Service) {
	avatarService = service
}

// ------------------------------------------------------------------------------------------------

func SetEmailService(service *email.Service) {
	emailService = service
}

// ------------------------------------------------------------------------------------------------
