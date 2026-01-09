package bootstrap

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"scifi-search/app/database"
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
	avatarService := getAvatarsService()
	emailService := getEmailService()

	initializeAuthorizationMechanisms()

	registerEndpoints(HTTPDependencies{
		Queries:       queries,
		AvatarService: avatarService,
		EmailService:  emailService,
	})

	loadLanguages()

	return &App{
		Resources: &Resources{
			DB:      db,
			Queries: queries,
		},
	}, nil
}

// ------------------------------------------------------------------------------------------------
