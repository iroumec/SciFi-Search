package supertokens

// ------------------------------------------------------------------------------------------------
// Imports
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"scifi-search/app/email"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// ------------------------------------------------------------------------------------------------
// Constants
// Constants
// ------------------------------------------------------------------------------------------------

const (
	WebsiteDomain = "http://localhost:8080"
)

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

var (
	PasswordEmptyError = errors.New("error.password-empty")
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

func Initialize(emailService *email.Service) {

	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: fmt.Sprintf("http://%s:%s",
				os.Getenv("SUPERTOKENS_HOST"),
				os.Getenv("SUPERTOKENS_PORT"),
			),
			APIKey: os.Getenv("SUPERTOKENS_API_KEY"),
		},

		AppInfo: supertokens.AppInfo{
			AppName:       "scifi-search",
			APIDomain:     WebsiteDomain,
			WebsiteDomain: WebsiteDomain,
		},

		RecipeList: []supertokens.Recipe{

			// The pressence of different roles is specified.
			userroles.Init(nil),

			// Email/Password sign-up and log-in is enabled.
			emailpassword.Init(&epmodels.TypeInput{
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "password",
							// Password validation.
							Validate: func(value interface{}, tenantId string) *string {
								password := value.(string)
								if len(password) < 1 {
									err := PasswordEmptyError.Error()
									return &err
								}
								return nil // Allows everything else.
							},
						},
					},
				},
			}),

			// Email verification configuration.
			emailverification.Init(evmodels.TypeInput{

				// This mode allows the user to use some functionalities of the page
				// despict not having its email verified.
				// If a mode restrictive mode is required, use `evmodels.ModeRequired`.
				Mode: evmodels.ModeOptional,

				EmailDelivery: &emaildelivery.TypeInput{
					Service: EmailDelivery(emailService),
				},
			}),

			// Session configuration.
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

// ------------------------------------------------------------------------------------------------
