package workers

// ------------------------------------------------------------------------------------------------

import (
	"log"
	"scifi-search/app/infra/email"
)

// ------------------------------------------------------------------------------------------------

// Un 'work' de envío de email.
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

// Canal global para la cola de emails.
// Se pueden encolar hasta N emails sin bloquear.
var emailQueue = make(chan EmailJob, numberOfMaxQueuedEmails)

// ------------------------------------------------------------------------------------------------

// Inicia el worker que procesa emails de forma asíncrona.
func StartEmailWorker() {
	go func() {
		for job := range emailQueue {
			// Se procesa cada email.
			err := email.Send(job.To, job.Subject, job.Body)
			if debug {
				if err != nil {
					log.Printf("Error enviando email a %s: %v", job.To, err)
					// TODO: implementar reintentos.
				} else {
					log.Printf("Email enviado exitosamente a %s", job.To)
				}
			}
		}
	}()
}

// ------------------------------------------------------------------------------------------------

// Encola un email para envío asíncrono.
func SendEmailAsync(to, subject, body string) {
	select {
	case emailQueue <- EmailJob{To: to, Subject: subject, Body: body}:
		// Email encolado exitosamente
	default:
		// La cola está llena.
		log.Printf("Cola de emails llena, no se pudo encolar email a %s", to)
	}
}

// ------------------------------------------------------------------------------------------------
