package handlers

// ---------------------------------------------------------------------
// Importaciones
// ---------------------------------------------------------------------

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"scifi-search/app/auth"
	sqlc "scifi-search/app/database"
	"scifi-search/app/http/cookies"
	"scifi-search/app/languages"
	"scifi-search/app/utils"
	"scifi-search/app/utils/checkers"
	"scifi-search/app/views"
	"scifi-search/app/workers"

	"github.com/a-h/templ"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// ---------------------------------------------------------------------
// Constantes
// ---------------------------------------------------------------------

func RegisterAuthenticationHandlers() {

	// Creación de la cuenta de administración.
	createAdmin()

	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/login", logInHandler)
	http.HandleFunc("/signout", signOutHandler)
	http.HandleFunc("/loader", signOutHandler)
	http.HandleFunc("/delete-account", deleteUserHandler)

	// Handler de verificación de email.
	http.HandleFunc("/auth/verify-email", verifyEmailHandler)

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
		component := views.SignUpPage("", auth.GetCurrentAuthorizationLevel(w, r, queries), languages.GetTranslatorFromRequest(r))
		templ.Handler(component).ServeHTTP(w, r)
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
	name := r.Form.Get("name")
	surname := r.Form.Get("surname")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Validación.
	if checkers.HayCampoIncompleto(name, surname, email, password) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	newUser, resp := createUser(name, surname, email, password, auth.UserRole.Name)

	if newUser != nil && resp != nil {

		// Creación de sesión en SuperTokens.
		_, err := session.CreateNewSession(r, w, "", resp.OK.User.ID, nil, nil)
		if err != nil {
			http.Error(w, "Error al crear la sesión: "+err.Error(), http.StatusInternalServerError)
			return
		}

		cookies.AddFlashCookie(w, languages.GetTranslatorFromRequest(r)("email-verification.sent"))
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------

func createAdmin() {

	// Obtención de los datos del administrador.
	name := utils.GetEnv("ADMIN_NAME", "admin")
	surname := utils.GetEnv("ADMIN_SURNAME", "full-access")
	email := utils.GetEnv("ADMIN_EMAIL", "admin@scifi-search.com")
	password := utils.GetEnv("ADMIN_PASSWORD", "admin")

	user, resp := createUser(name, surname, email, password, auth.AdminRole.Name)

	if user == nil || resp == nil {
		log.Fatal("Ocurrió un error al momento de crear al usuario.")
	}
}

// ------------------------------------------------------------------------------------------------

func createUser(name, surname, email, password, role string) (*sqlc.User, *epmodels.SignUpResponse) {

	var newUser sqlc.User

	resp, err := emailpassword.SignUp("", email, password)
	if err != nil {
		log.Fatal("It wasn't possible to create the user.")
	}

	// Verificación de email ya registrado.
	if resp.EmailAlreadyExistsError != nil {
		log.Print("Email already in use.")
		return &newUser, &resp
	}

	if resp.OK != nil {
		_, err = queries.CreateUser(context.TODO(), sqlc.CreateUserParams{
			Name:    name,
			Surname: surname,
			AuthID:  resp.OK.User.ID,
		})
		if err != nil {
			log.Println("An unexpected error ocurred during the creation of the administrator's account:", err)
		}

		// Creación del token de verificación.
		auth.SendVerificationEmail(resp.OK.User.ID, email)
	}

	// Asignación del rol al usuario.
	userroles.AddRoleToUser("public", resp.OK.User.ID, role, nil)

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
		cookies.AddFlashCookie(w, "¡Email verificado exitosamente!")
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
	translator := languages.GetTranslatorFromRequest(r)

	if r.Method == http.MethodGet {
		views.LoginPage(auth.GetCurrentAuthorizationLevel(w, r, queries), translator).Render(r.Context(), w)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Se parsean y obtienen las credenciales.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Se intenta realizar un log in.
	resp, err := emailpassword.SignIn("", email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Credenciales incorrectas.
	if resp.WrongCredentialsError != nil {
		cookies.AddFlashCookie(w, "Credenciales incorrectas")
		views.LoginPage(auth.GetCurrentAuthorizationLevel(w, r, queries), translator).Render(r.Context(), w)
		return
	}

	// Login exitoso: se crea la sesión y se redirige.
	if resp.OK != nil {
		_, err := session.CreateNewSession(r, w, "", resp.OK.User.ID, nil, nil)
		if err != nil {
			http.Error(w, "Error al crear la sesión", http.StatusInternalServerError)
			return
		}

		cookies.AddFlashCookie(w, "Welcome back!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Error desconocido", http.StatusInternalServerError)
}

// ------------------------------------------------------------------------------------------------
// Sign Out (Log Out)
// ------------------------------------------------------------------------------------------------

func signOutHandler(w http.ResponseWriter, r *http.Request) {

	auth.RevokeSession(w, r)

	cookies.AddFlashCookie(w, "Successful signout!")

	http.Redirect(w, r, "/", http.StatusFound)
}

// ------------------------------------------------------------------------------------------------
// Get Current User
// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deleteUser(w, r)
}

// ------------------------------------------------------------------------------------------------

func deleteUser(w http.ResponseWriter, r *http.Request) {

	user, err := auth.GetCurrentUser(w, r, queries)
	if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	userEmail := auth.GetCurrentUserEmail(w, r)
	if userEmail == nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Se elimina el usuario de supertokens.
	err = deleteSupertokensUser(user.AuthID)
	if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Se elimina el avatar del usuario (si tiene).
	err = avatarService.Delete(r.Context(), user.UserID)
	if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	// Se elimina el usuario de la base de datos.
	err = queries.DeleteUser(r.Context(), user.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: El usuario no existe.
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			// Error 500: Hubo un problema con la base de datos u otro error inesperado.
			log.Printf("Error al obtener usuario por ID %d: %v", user.UserID, err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	// Se cierra la sesión del usuario.
	auth.RevokeSession(w, r)

	cookies.AddFlashCookie(w, "Usuario eliminado. ¡Lamentamos que te vayas!")

	workers.SendEmailAsync(
		*userEmail,
		"¡Lamentamos que te vayas!",
		"Nos entristece ver que te vayas. Para tu seguridad, hemos eliminado todos tus datos. ¡Esperamos volver a verte pronto!",
	)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------

func deleteSupertokensUser(authID string) error {
	err := supertokens.DeleteUser(authID)

	if err != nil {
		return err
	}

	// Usuario eliminado exitosamente.
	return nil
}

// ------------------------------------------------------------------------------------------------
