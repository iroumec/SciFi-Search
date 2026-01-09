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
	handlers.SetAvatarService(dependencies.AvatarService)
	handlers.SetEmailService(dependencies.EmailService)

	handlers.RegisterStatic()
	handlers.RegisterIndex()
	handlers.RegisterHealth()
	handlers.RegisterAuthenticationHandlers()
	handlers.RegisterSearchHandlers()
	handlers.RegisterSettingsHandlers()
	handlers.RegisterAvatarHandlers()
	handlers.RegisterHistoryHandlers()
	handlers.RegisterFundingHandlers()
	handlers.RegisterTrendsHandlers()
	handlers.RegisterLanguageHandlers()
}

// ------------------------------------------------------------------------------------------------
