package bootstrap

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"scifi-search/app/database"
	"scifi-search/app/languages"
)

// ------------------------------------------------------------------------------------------------
// Estructuras
// ------------------------------------------------------------------------------------------------

type App struct {
	Resources *Resources
}

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

func Boot() (*App, error) {

	db, err := initializeDatabase()
	if err != nil {
		return nil, err
	}

	queries := database.New(db)
	emailService := getEmailService()
	avatarService := getAvatarsService()

	initializeAuthorizationMechanisms()

	registerEndpoints(HTTPDependencies{
		Queries:       queries,
		AvatarService: avatarService,
		EmailService:  emailService,
	})

	languages.LoadAllMessages()

	return &App{
		Resources: &Resources{
			DB:      db,
			Queries: queries,
		},
	}, nil
}

// ------------------------------------------------------------------------------------------------
