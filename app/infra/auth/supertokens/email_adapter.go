package supertokens

import (
	"scifi-search/app/email"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func EmailDelivery(service *email.Service) *emaildelivery.EmailDeliveryInterface {

	send := func(input emaildelivery.EmailType, _ supertokens.UserContext) error {
		if input.EmailVerification != nil {
			ev := input.EmailVerification
			return service.SendVerification(
				ev.User.Email,
				ev.EmailVerifyLink,
			)
		}
		return nil
	}

	return &emaildelivery.EmailDeliveryInterface{
		SendEmail: &send,
	}
}
