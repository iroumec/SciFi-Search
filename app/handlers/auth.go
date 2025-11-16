package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	emailpassword "github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"

	sqlc "tpe/web/app/database"
	"tpe/web/app/views"

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
			session.Init(nil),
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

// ---------------------------------------------------------------------

func registerAuthenticationHandlers() {

	initializeSupertokens()

	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/signin", signIn)

	http.HandleFunc("/auth/session/refresh", func(w http.ResponseWriter, r *http.Request) {
		session.RefreshSession(r, w)
	})

	http.HandleFunc("/auth/sessioninfo", func(w http.ResponseWriter, r *http.Request) {
		sessionContainer := session.GetSessionFromRequestContext(r.Context())
		if sessionContainer == nil {
			http.Error(w, "No session", http.StatusUnauthorized)
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
	})
}

// ---------------------------------------------------------------------

func signUp(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		component := views.SignUpPage("")
		templ.Handler(component).ServeHTTP(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Surname  string `json:"surname"`
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

	if resp.OK != nil {
		_, err := queries.CreateUser(r.Context(), sqlc.CreateUserParams{
			Name:    body.Name,
			Surname: body.Surname,
		})
		if err != nil {
			log.Println("Error creando usuario interno:", err)
		}
	}

	json.NewEncoder(w).Encode(resp)
}

// ---------------------------------------------------------------------

func signIn(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		component := views.LoginPage("")
		templ.Handler(component).ServeHTTP(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

		// Si es una API, responde con éxito
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"userId": userID,
		})
		return
	}

	// Manejo de cualquier otro caso inesperado.
	http.Error(w, "Error desconocido durante el inicio de sesión", http.StatusInternalServerError)
}
