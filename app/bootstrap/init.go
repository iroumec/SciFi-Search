package bootstrap

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"scifi-search/app/database"
	"scifi-search/app/languages"
)

// ------------------------------------------------------------------------------------------------
// Structures
// ------------------------------------------------------------------------------------------------

type App struct {
	Resources *Resources
}

// ------------------------------------------------------------------------------------------------
// Services
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

	languages.LoadAllMessagesFromFolders()

	registerEndpoints(HTTPDependencies{
		Queries:       queries,
		AvatarService: avatarService,
		EmailService:  emailService,
	})

	return &App{
		Resources: &Resources{
			DB:      db,
			Queries: queries,
		},
	}, nil
}

// ------------------------------------------------------------------------------------------------
