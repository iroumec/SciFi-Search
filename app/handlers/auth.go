package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	emailpassword "github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"

	"scifi-search/app/database"
	sqlc "scifi-search/app/database"
	"scifi-search/app/utils"
	"scifi-search/app/utils/email"
	"scifi-search/app/views"

	"github.com/a-h/templ"
)

const (
	websiteDomain = "http://localhost:8080"
)

// ---------------------------------------------------------------------

func initializeSupertokens() {

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

			// Se permite inicio de sesión mediante email/password.
			emailpassword.Init(nil),

			// Se permite inicio de sesión mediante cuenta de terceros.
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.ProviderInput{
						{
							Config: tpmodels.ProviderConfig{
								ThirdPartyId: "google",
								Clients: []tpmodels.ProviderClientConfig{
									{
										ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
										ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
									},
								},
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

	if err != nil {
		log.Fatal(err)
	}
}

// ---------------------------------------------------------------------

func registerAuthenticationHandlers() {

	initializeSupertokens()

	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/login", logInHandler)
	http.HandleFunc("/signout", signOutHandler)

	// Handler de verificación de email.
	http.HandleFunc("/auth/verify-email", verifyEmailHandler)

	// Handler para Google OAuth callback.
	http.HandleFunc("/auth/callback/google", googleCallbackHandler)

	http.HandleFunc("/auth/session/refresh", func(w http.ResponseWriter, r *http.Request) {
		session.RefreshSession(r, w)
	})

	// Al usar VerifySession, la sesión ya está garantizada y puesta en el contexto
	// si las cookies de sesión son válidas.
	http.HandleFunc("/auth/sessioninfo", session.VerifySession(nil, func(w http.ResponseWriter, r *http.Request) {
		// Al usar session.VerifySession, la sesión ya estará garantizada y
		// puesta en el contexto si las cookies de sesión son válidas.

		sessionContainer := session.GetSessionFromRequestContext(r.Context())

		if sessionContainer == nil {
			// Este caso sólo debería ocurrir si hay un error interno en el middleware
			// o si VerifySession no está configurado para devolver error 401.
			http.Error(w, "No session (Error interno o configuración)", http.StatusUnauthorized)
			return
		}

		userID := sessionContainer.GetUserID()
		rawPayload := sessionContainer.GetAccessTokenPayload()

		payload := make(map[string]string)
		for k, v := range rawPayload {
			payload[k] = fmt.Sprintf("%v", v)
		}

		component := views.SessionInfo(userID, payload)
		templ.Handler(component).ServeHTTP(w, r)
	}))
}

// ------------------------------------------------------------------------------------------------
// Sign Up (Registro)
// ------------------------------------------------------------------------------------------------

func signUpHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		component := views.SignUpPage("", utils.GetTranslatorFromRequest(r))
		templ.Handler(component).ServeHTTP(w, r)
		return
	}

	newUser, resp := createUser(w, r)

	if newUser != nil && resp != nil {

		// Creación de sesión en SuperTokens.
		_, err := session.CreateNewSession(r, w, "", resp.OK.User.ID, nil, nil)
		if err != nil {
			http.Error(w, "Error al crear la sesión: "+err.Error(), http.StatusInternalServerError)
			return
		}

		utils.AddFlashCookie(w, utils.GetTranslatorFromRequest(r)("email-verification.sent"))
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------

func createUser(w http.ResponseWriter, r *http.Request) (*sqlc.User, *epmodels.SignUpResponse) {

	var newUser sqlc.User

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return &newUser, nil
	}

	// Parseo del formulario enviado por POST.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario: "+err.Error(), http.StatusBadRequest)
		return &newUser, nil
	}

	// Obtención de los datos del usuario.
	name := r.Form.Get("name")
	surname := r.Form.Get("surname")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Validación.
	if utils.HayCampoIncompleto(name, surname, email, password) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return &newUser, nil
	}

	resp, err := emailpassword.SignUp("", email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return &newUser, nil
	}

	// Verificación de email ya registrado.
	if resp.EmailAlreadyExistsError != nil {
		utils.AddFlashCookie(w, utils.GetTranslatorFromRequest(r)("Usuario ya registado en el sistema. Inicie sesión."))
		return &newUser, nil
	}

	if resp.OK != nil {
		newUser, err = queries.CreateUser(r.Context(), sqlc.CreateUserParams{
			Name:    name,
			Surname: surname,
			AuthID:  resp.OK.User.ID,
		})
		if err != nil {
			log.Println("Error creando usuario interno:", err)
		}

		// Se crea token de verificación.
		tokenResponse, err := emailverification.CreateEmailVerificationToken("", resp.OK.User.ID, &email, nil)
		if err != nil {
			log.Println("Error creando token de verificación:", err)
		} else if tokenResponse.OK != nil {
			// Se envía email de verificación.
			err = emailverification.SendEmail(emaildelivery.EmailType{
				EmailVerification: &emaildelivery.EmailVerificationType{
					User: emaildelivery.User{
						ID:    resp.OK.User.ID,
						Email: email,
					},
					EmailVerifyLink: websiteDomain + "/auth/verify-email?token=" + tokenResponse.OK.Token,
				},
			}, nil)

			if err != nil {
				log.Println("Error enviando email de verificación:", err)
			}
		}
	}

	return &newUser, &resp
}

// ------------------------------------------------------------------------------------------------

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "Token no proporcionado", http.StatusBadRequest)
		return
	}

	// Se verifica el token.
	response, err := emailverification.VerifyEmailUsingToken("", token, nil)
	if err != nil {
		log.Println("Error verificando email:", err)
		http.Error(w, "Error verificando email", http.StatusInternalServerError)
		return
	}

	if response.OK != nil {
		utils.AddFlashCookie(w, "¡Email verificado exitosamente!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if response.EmailVerificationInvalidTokenError != nil {
		http.Error(w, "Token inválido o expirado", http.StatusBadRequest)
	} else {
		http.Error(w, "Error desconocido al verificar email", http.StatusInternalServerError)
	}
}

// ------------------------------------------------------------------------------------------------
// Sign In (Log In)
// ------------------------------------------------------------------------------------------------

func logInHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		component := views.LoginPage("", utils.GetTranslatorFromRequest(r))
		templ.Handler(component).ServeHTTP(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parseo del formulario enviado por POST.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Obtención de los datos del usuario.
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Inicio de sesión.
	resp, err := emailpassword.SignIn("", email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Abarca tanto la comprovación del email como de la constraseña.
	if resp.WrongCredentialsError != nil {
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	// Si todo está OK, se crea la sesión.
	if resp.OK != nil {
		userID := resp.OK.User.ID

		// Se crea la sesión y automáticamente se establecen
		// las cookies de sesión en el http.ResponseWriter (w).
		_, err := session.CreateNewSession(r, w, "", userID, nil, nil)
		if err != nil {
			http.Error(w, "Error al crear la sesión: "+err.Error(), http.StatusInternalServerError)
			return
		}

		utils.AddFlashCookie(w, "Welcome back!")

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Manejo de cualquier otro caso inesperado.
	http.Error(w, "Error desconocido durante el inicio de sesión", http.StatusInternalServerError)
}

// ------------------------------------------------------------------------------------------------
// Sign Out (Log Out)
// ------------------------------------------------------------------------------------------------

func signOutHandler(w http.ResponseWriter, r *http.Request) {
	sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
		SessionRequired: boolPtr(false), // False -> No error si no hay sessión. Posiblemente deba cambiarse luego.
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if sessionContainer == nil {
		http.Error(w, "No hay sesión activa", http.StatusUnauthorized)
		return
	}

	if err := sessionContainer.RevokeSession(); err != nil {
		http.Error(w, "Error cerrando sesión", http.StatusInternalServerError)
		return
	}

	utils.AddFlashCookie(w, "Successful signout!")

	http.Redirect(w, r, "/", http.StatusFound)
}

// ------------------------------------------------------------------------------------------------
// Funciones Auxiliares
// ------------------------------------------------------------------------------------------------

// Retorna si el usuario está autenticado.
func isUserAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	// Intentar obtener la sesión sin requerirla
	sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
		SessionRequired: boolPtr(false),
	})

	return err == nil && sessionContainer != nil
}

// ---------------------------------------------------------------------

// Necesaria para obtener un puntero a un booleano,
// el cual Supertokens requiere.
func boolPtr(b bool) *bool {
	return &b
}

// ------------------------------------------------------------------------------------------------
// Get Current User
// ------------------------------------------------------------------------------------------------

func getCurrentUser(w http.ResponseWriter, r *http.Request) *database.User {

	if isUserAuthenticated(w, r) {

		sessionContainer, _ := session.GetSession(r, nil, &sessmodels.VerifySessionOptions{
			SessionRequired: boolPtr(false),
		})

		supertokensUserID := sessionContainer.GetUserID()

		user, err := queries.GetUserByAuthID(r.Context(), supertokensUserID)
		if err != nil {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}

		return &user

	} else {

		return nil
	}
}

// ------------------------------------------------------------------------------------------------

func deleteUser(w http.ResponseWriter, r *http.Request, id int32) {

	err := queries.DeleteUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: El usuario no existe.
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			// Error 500: Hubo un problema con la base de datos u otro error inesperado.
			log.Printf("Error al obtener usuario por ID %d: %v", id, err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Se eliminó al usuario de ID %d.", id)

	// Por defecto, la respuesta es 200 OK.
}

// ------------------------------------------------------------------------------------------------

// Función auxiliar para verificar si el email está verificado.
func isEmailVerified(w http.ResponseWriter, r *http.Request) bool {
	sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
		SessionRequired: boolPtr(true),
	})
	if err != nil {
		return false
	}

	if sessionContainer == nil {
		return false
	}

	userID := sessionContainer.GetUserID()
	isVerified, err := emailverification.IsEmailVerified(userID, nil, nil)

	return isVerified
}

// ------------------------------------------------------------------------------------------------

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	verified := isEmailVerified(w, r)
	if !verified {
		http.Error(w, "Debes verificar tu email para acceder a esta funcionalidad", http.StatusForbidden)
		return
	}
}

// ------------------------------------------------------------------------------------------------

func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {

	// SuperTokens maneja esto automáticamente si se usa su frontend.
	// Si se implementa un frontend propio, se necesita procesar el callback aquí.

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Código no proporcionado", http.StatusBadRequest)
		return
	}

	// Redirigir a la página principal.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getUserEmail(userID string) *string {

	user, err := emailpassword.GetUserByID(userID)
	if err != nil || user == nil {
		return nil
	}

	return &user.Email
}
