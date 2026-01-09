package handlers

import (
	"scifi-search/app/avatars"
	sqlc "scifi-search/app/database"
	"scifi-search/app/email"
)

// ------------------------------------------------------------------------------------------------
// Constantes
// ------------------------------------------------------------------------------------------------

// Ruta a partir de la cual se servirán los archivos estáticos.
const (
	debug   = false
	fileDir = "./static"
)

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

var (
	queries       *sqlc.Queries
	avatarService *avatars.Service
	emailService  *email.Service
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
