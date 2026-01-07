package bootstrap

import (
	"scifi-search/app/database"
)

type App struct {
	Resources *Resources
}

func Boot() (*App, error) {

	db, err := initializeDatabase()
	if err != nil {
		return nil, err
	}

	queries := database.New(db)
	avatarService := getAvatarsService()
	emailService := getEmailService()

	initializeAuthorizationMechanisms()
	startWorkers(emailService)

	registerEndpoints(HTTPDependencies{
		Queries:       queries,
		AvatarService: avatarService,
	})

	return &App{
		Resources: &Resources{
			DB:      db,
			Queries: queries,
		},
	}, nil
}
