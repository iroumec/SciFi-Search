package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"scifi-search/app/infra/email"
	"slices"

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

const (
	websiteDomain = "http://localhost:8080"
)

// Role representa un rol del sistema con su nombre y nivel de autorización.
type Role struct {
	Name  string
	Level int
}

// Roles predefinidos del sistema.
var (
	NoRole     = Role{Name: "no-role", Level: -1}
	AdminRole  = Role{Name: "admin", Level: 2}
	LoaderRole = Role{Name: "loader", Level: 1}
	UserRole   = Role{Name: "user", Level: 0}
)

func InitializeSupertokens() {

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
			APIDomain:     websiteDomain,
			WebsiteDomain: websiteDomain,
		},

		RecipeList: []supertokens.Recipe{

			// Se especifica la utilización de distintos roles de usuarios.
			userroles.Init(nil),

			// Se permite inicio de sesión mediante email/password.
			emailpassword.Init(&epmodels.TypeInput{
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "password",
							// Validación muy permisiva para probar rápido.
							Validate: func(value interface{}, tenantId string) *string {
								password := value.(string)
								if len(password) < 1 {
									err := "La contraseña no puede estar vacía"
									return &err
								}
								return nil // Se acepta todo lo demás.
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
					Service: email.NewMailHogService(),
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

	// Creación de roles.
	userroles.CreateNewRoleOrAddPermissions(AdminRole.Name, []string{"full-access"}, nil)            // Administrador.
	userroles.CreateNewRoleOrAddPermissions(LoaderRole.Name, []string{"manage-own-financings"}, nil) // Entidades que cargan financiamiento.
	userroles.CreateNewRoleOrAddPermissions(UserRole.Name, []string{"view-only"}, nil)               // Usuarios normales.

	if err != nil {
		log.Fatal(err)
	}
}

func GetAuthenticationLevel(userID string) int {
	roles, _ := userroles.GetRolesForUser("public", userID, nil)

	if slices.Contains(roles.OK.Roles, AdminRole.Name) {
		return AdminRole.Level
	} else if slices.Contains(roles.OK.Roles, LoaderRole.Name) {
		return LoaderRole.Level
	} else if slices.Contains(roles.OK.Roles, UserRole.Name) {
		return UserRole.Level
	}

	// Usuario sin autenticar.
	return NoRole.Level
}
