package supertokens

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
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
// ------------------------------------------------------------------------------------------------

const (
	WebsiteDomain = "http://localhost:8080"
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

			// Specifying the use of different user roles.
			userroles.Init(nil),

			// Allowing log-in through email/password
			emailpassword.Init(&epmodels.TypeInput{
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "password",
							// Quick and small validation
							Validate: func(value interface{}, tenantId string) *string {
								password := value.(string)
								if len(password) < 1 {
									err := "The password cannot be empty." //TODO: como se usaria el traductor aca?
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
				// El siguiente modo permite que el usuario use la página normalmente
				// aunque no haya verificado su email, aunque se le pueden prohibir ciertas
				// accioens manualmente.
				Mode: evmodels.ModeOptional,

				// El siguiente estado inhabilita la sesión a menos de que el email se
				// encuentre verificado. Es el más estricto.
				// Mode: evmodels.ModeRequired,

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
