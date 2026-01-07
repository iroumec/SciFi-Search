package mailhog

import (
	"net/smtp"
	"os"

	"scifi-search/app/email"
)

type Provider struct{}

func New() email.Provider {
	return &Provider{}
}

func (p *Provider) Send(to, subject, body string) error {
	msg := []byte(
		"From: noreply@scifi-search.com\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n\r\n" +
			body,
	)

	return smtp.SendMail(
		os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"),
		nil,
		"noreply@scifi-search.com",
		[]string{to},
		msg,
	)
}
