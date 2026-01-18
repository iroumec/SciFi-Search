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
	MethodNotAllowedError        = errors.New("errors.method-not-allowed")
	FormParsingError             = errors.New("errors.form-parsing")
	UnexpectedError              = errors.New("errors.unexpected")
	UnknownError                 = errors.New("errors.unknown")
	InternalServerError          = errors.New("errors.internal-server")
	RequiredDataNotSpecified     = errors.New("errors.required-data-not-specified")
	MissingRequiredFieldsError   = errors.New("errors.missing-required-fields")
	UnsupportedExportFormatError = errors.New("errors.unsupported-export-format")
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
