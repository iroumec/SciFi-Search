package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	emailpassword "github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"

	sqlc "scifi-search/app/database"
	"scifi-search/app/utils"
	"scifi-search/app/views"

	"github.com/a-h/templ"
)

// ---------------------------------------------------------------------

func initializeSupertokens() {

	websiteDomain := "http://localhost:8080"

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
			emailpassword.Init(nil),
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
	http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/signout", signOutHandler)
	http.HandleFunc("/create-user", userCreationHandler)

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
		component := views.SignUpPage("")
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

		component := views.SuccessfulSignUpPage()
		templ.Handler(component).ServeHTTP(w, r)
	} else {

		// Manejo de cualquier otro caso inesperado.
		http.Error(w, "Error desconocido durante el registro.", http.StatusInternalServerError)
	}
}

// ------------------------------------------------------------------------------------------------

func userCreationHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		component := views.SignUpPage("")
		templ.Handler(component).ServeHTTP(w, r)
		return
	}

	newUser, _ := createUser(w, r)
	if newUser != nil {
		component := views.UserIndividual(*newUser)
		templ.Handler(component).ServeHTTP(w, r)
	} else {
		// Manejo de cualquier otro caso inesperado.
		http.Error(w, "Error desconocido durante la creación del usuario.", http.StatusInternalServerError)
	}
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

	if resp.OK != nil {
		newUser, err = queries.CreateUser(r.Context(), sqlc.CreateUserParams{
			Name:    name,
			Surname: surname,
			AuthID:  resp.OK.User.ID,
		})
		if err != nil {
			log.Println("Error creando usuario interno:", err)
		}
	}

	return &newUser, &resp
}

// ------------------------------------------------------------------------------------------------
// Sign In (Log In)
// ------------------------------------------------------------------------------------------------

func signInHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		component := views.LoginPage("")
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

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
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

	http.Redirect(w, r, "/", http.StatusFound)
}

// ------------------------------------------------------------------------------------------------
// Funciones Auxiliares
// ------------------------------------------------------------------------------------------------

// Retorna si el usuario está autenticado.
func isUserAuthenticated(r *http.Request) bool {
	// Intentar obtener la sesión sin requerirla
	sessionContainer, err := session.GetSession(r, nil, &sessmodels.VerifySessionOptions{
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

// ---------------------------------------------------------------------
