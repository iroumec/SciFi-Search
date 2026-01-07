package bootstrap

import (
	"scifi-search/app/email"
	"scifi-search/app/workers"
)

func startWorkers(emailService *email.Service) {

	workers.SetEmailService(emailService)
	// Se inicializa el worker de envíos asíncronos de emails.
	workers.StartEmailWorker()
}
