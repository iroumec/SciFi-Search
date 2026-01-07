package workers

import (
	"log"

	"scifi-search/app/email"
)

// ------------------------------------------------------------------------------------------------

type EmailJob struct {
	To      string
	Subject string
	Body    string
}

// ------------------------------------------------------------------------------------------------

const (
	debug                   = true
	numberOfMaxQueuedEmails = 100
)

// ------------------------------------------------------------------------------------------------

var (
	emailQueue   = make(chan EmailJob, numberOfMaxQueuedEmails)
	emailService *email.Service
)

// ------------------------------------------------------------------------------------------------

// Inyección del servicio
func SetEmailService(s *email.Service) {
	emailService = s
}

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

func SendEmailAsync(to, subject, body string) {
	select {
	case emailQueue <- EmailJob{To: to, Subject: subject, Body: body}:
	default:
		log.Printf("Cola de emails llena, no se pudo encolar email a %s", to)
	}
}
