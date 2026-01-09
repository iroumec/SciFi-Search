package bootstrap

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"scifi-search/app/email"
	"scifi-search/app/infra/email/mailhog"
	"scifi-search/app/workers"
)

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

func getEmailService() *email.Service {

	jobQueue := make(chan workers.Job, 100)

	workers.StartWorker(jobQueue)

	provider := mailhog.New()
	asyncProvider := email.NewAsyncProvider(jobQueue, provider)

	return email.New(asyncProvider)
}

// ------------------------------------------------------------------------------------------------
