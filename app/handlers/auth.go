package handlers

// ------------------------------------------------------------------------------------------------
// TODO: a desarrollar en etapas posteriores.
// ------------------------------------------------------------------------------------------------

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

    "github.com/supertokens/supertokens-golang/supertokens"
    "github.com/supertokens/supertokens-golang/recipe/emailpassword"
    "github.com/supertokens/supertokens-golang/recipe/session"

	//sqlc "tpe/web/app/database"

	"tpe/web/app/views"
	"github.com/a-h/templ"
)

// ------------------------------------------------------------------------------------------------

func initializeSupertokens() {
    websiteBaseURL := "http://localhost:8080"

    err := supertokens.Init(supertokens.TypeInput{
        Supertokens: &supertokens.ConnectionInfo{
            ConnectionURI: fmt.Sprintf("http://%s:%s",
                os.Getenv("SUPERTOKENS_HOST"),
                os.Getenv("SUPERTOKENS_PORT"),
            ),
            APIKey: os.Getenv("SUPERTOKENS_API_KEY"),
        },

        AppInfo: supertokens.AppInfo{
            AppName:         "scifi-search",
            APIDomain:       websiteBaseURL,
            WebsiteBasePath: &websiteBaseURL,
        },

        RecipeList: []supertokens.Recipe{
            emailpassword.Init(nil),/*emailpassword.Init(&emailpassword.TypeInput{
                Override: &emailpassword.OverrideStruct{
                    APIs: overrideSignUp,
                },
            }),*/
            session.Init(nil),
        },
    })

    if err != nil {
        log.Fatal(err)
    }
}

// ------------------------------------------------------------------------------------------
// OVERRIDE DEL SIGNUP
// ------------------------------------------------------------------------------------------

/*func overrideSignUp(original emailpassword.APIInterface) emailpassword.APIInterface {
    original.SignUpPOST = func(
        w http.ResponseWriter,
        r *http.Request,
        options emailpassword.APIOptions,
        userContext supertokens.UserContext,
    ) (*emailpassword.SignUpPOSTResponse, error) {

        resp, err := emailpassword.SignUpPOST(w, r, options, userContext)
        if err != nil {
            return resp, err
        }

        if resp.User != nil {
            handleInternalUserCreation(r, resp.User.ID, resp.User.Email)
        }

        return resp, nil
    }

    return original
}


// ------------------------------------------------------------------------------------------
// Creación Interna del Usuario (Nombre y Apellido)
// ------------------------------------------------------------------------------------------

func handleInternalUserCreation(r *http.Request, userID string, email string) {
    name := r.FormValue("name")
    surname := r.FormValue("surname")

    _, err := queries.CreateUser(r.Context(), sqlc.CreateUserParams{
        Name:    name,
        Surname: surname,
    })

    if err != nil {
        log.Println("Error creando usuario interno:", err)
    }
}*/

// ------------------------------------------------------------------------------------------------

func initSupertokens() {

	initializeSupertokens()

	http.HandleFunc("/signup", signUp)

	http.HandleFunc("/signin", signIn)

	http.HandleFunc("/auth/session/refresh", func(w http.ResponseWriter, r *http.Request) {
		session.RefreshSession(r, w)
	})

	http.HandleFunc("/auth/sessioninfo", func(w http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		if sessionContainer == nil {
			//http.Error(w, "No session", http.StatusUnauthorized)
			component := views.NoSessionInfo()
			templ.Handler(component).ServeHTTP(w, r)
			return
		}

		userID := sessionContainer.GetUserID()
		rawPayload := sessionContainer.GetAccessTokenPayload()

		// Conversión segura a map[string]string
		payload := make(map[string]string)
		for k, v := range rawPayload {
			payload[k] = fmt.Sprintf("%v", v)
		}

		// Render del templ
		component := views.SessionInfo(userID, payload)
		templ.Handler(component).ServeHTTP(w, r)
	})
}

// ------------------------------------------------------------------------------------------------

func signUp(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Leer datos del body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := emailpassword.SignUp("", body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

// ------------------------------------------------------------------------------------------------

func signIn(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Leer datos del body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := emailpassword.SignIn("", body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}
