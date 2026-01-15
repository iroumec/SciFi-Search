package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"errors"
	"scifi-search/app/avatars"
	sqlc "scifi-search/app/database"
	"scifi-search/app/email"
)

// ------------------------------------------------------------------------------------------------
// Constants
// ------------------------------------------------------------------------------------------------

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

// Errors.
var (
	MethodNotAllowedError = errors.New("error.method-not-allowed")
	FormParsingError      = errors.New("error.form-parsing")
	UnexpectedError       = errors.New("error.unexpected")
	UnknownError          = errors.New("error.unknown")
	InternalServerError   = errors.New("error.internal-server")
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
