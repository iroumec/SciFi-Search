package bootstrap

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"scifi-search/app/avatars"
	sqlc "scifi-search/app/database"
	"scifi-search/app/email"
	"scifi-search/app/http/handlers"
)

// ------------------------------------------------------------------------------------------------
// Estructuras
// ------------------------------------------------------------------------------------------------

type HTTPDependencies struct {
	Queries       *sqlc.Queries
	AvatarService *avatars.Service
	EmailService  *email.Service
}

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

func registerEndpoints(dependencies HTTPDependencies) {

	handlers.SetQueries(dependencies.Queries)
	handlers.SetEmailService(dependencies.EmailService)
	handlers.SetAvatarService(dependencies.AvatarService)

	handlers.RegisterIndex()
	handlers.RegisterStatic()
	handlers.RegisterHealth()
	handlers.RegisterSearchHandlers()
	handlers.RegisterTrendsHandlers()
	handlers.RegisterHistoryHandlers()
	handlers.RegisterFundingHandlers()
	handlers.RegisterSettingsHandlers()
	handlers.RegisterLanguageHandlers()
	handlers.RegisterExportationHandlers()
	handlers.RegisterAuthenticationHandlers()
}

// ------------------------------------------------------------------------------------------------
