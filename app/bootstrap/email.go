package bootstrap

import (
	"scifi-search/app/email"
	"scifi-search/app/infra/email/mailhog"
)

func getEmailService() *email.Service {

	provider := mailhog.New()

	return email.New(provider)
}
