package email

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"net/smtp"
	"os"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// ------------------------------------------------------------------------------------------------
// Constantes
// ------------------------------------------------------------------------------------------------

const (
	fromEmail = "noreply@scifi-search.com"
)

// ------------------------------------------------------------------------------------------------
// Estructuras
// ------------------------------------------------------------------------------------------------

type MailHogService struct{}

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

func Send(to, subject, body string) error {

	msg := []byte(
		"From: " + fromEmail + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
			body,
	)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	return smtp.SendMail(
		smtpHost+":"+smtpPort, // MailHog SMTP
		nil,                   // Sin autorización.
		fromEmail,
		[]string{to},
		msg,
	)
}

// ------------------------------------------------------------------------------------------------

// Pedido de esta forma por el SDK de MailHog.
func NewMailHogService() *emaildelivery.EmailDeliveryInterface {
	send := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.EmailVerification != nil {
			ev := input.EmailVerification
			return Send(
				ev.User.Email,
				"Verificá tu email",
				"Hacé click acá:\n\n"+ev.EmailVerifyLink,
			)
		}
		return nil
	}

	return &emaildelivery.EmailDeliveryInterface{
		SendEmail: &send,
	}
}

// ------------------------------------------------------------------------------------------------
