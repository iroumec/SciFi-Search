package workers

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"log"

	"scifi-search/app/email"
)

// ------------------------------------------------------------------------------------------------
// Estructuras
// ------------------------------------------------------------------------------------------------

type EmailJob struct {
	To      string
	Subject string
	Body    string
}

// ------------------------------------------------------------------------------------------------
// Constantes
// ------------------------------------------------------------------------------------------------

const (
	debug                   = true
	numberOfMaxQueuedEmails = 100
)

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

var (
	emailQueue   = make(chan EmailJob, numberOfMaxQueuedEmails)
	emailService *email.Service
)

// ------------------------------------------------------------------------------------------------
// Setters
// ------------------------------------------------------------------------------------------------

// Inyección del servicio
func SetEmailService(service *email.Service) {
	emailService = service
}

// ------------------------------------------------------------------------------------------------
// Servicios
// ------------------------------------------------------------------------------------------------

// Inicia el worker
func StartEmailWorker() {
	go func() {
		for job := range emailQueue {

			if emailService == nil {
				log.Println("email service not configured")
				continue
			}

			err := emailService.Send(job.To, job.Subject, job.Body)

			if debug {
				if err != nil {
					log.Printf("Error enviando email a %s: %v", job.To, err)
				} else {
					log.Printf("Email enviado exitosamente a %s", job.To)
				}
			}
		}
	}()
}

// ------------------------------------------------------------------------------------------------

// Envía un email asíncrono.
func SendEmailAsync(to, subject, body string) {
	select {
	case emailQueue <- EmailJob{To: to, Subject: subject, Body: body}:
	default:
		log.Printf("Cola de emails llena, no se pudo encolar email a %s", to)
	}
}

// ------------------------------------------------------------------------------------------------
