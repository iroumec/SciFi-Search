package handlers

import (
	"scifi-search/app/avatars"
	sqlc "scifi-search/app/database"
	"scifi-search/app/email"
)

// ------------------------------------------------------------------------------------------------
// Constantes del paquete
// ------------------------------------------------------------------------------------------------

// Ruta a partir de la cual se servirán los archivos estáticos.
const (
	debug   = false
	fileDir = "./static"
)

var (
	queries       *sqlc.Queries
	avatarService *avatars.Service
	emailService  *email.Service
)

func SetQueries(q *sqlc.Queries) {
	queries = q
}

func SetAvatarService(service *avatars.Service) {
	avatarService = service
}

func SetEmailService(service *email.Service) {
	emailService = service
}
