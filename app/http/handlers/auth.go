package handlers

// ---------------------------------------------------------------------
// Importaciones
// ---------------------------------------------------------------------

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"scifi-search/app/auth"
	"scifi-search/app/database"
	sqlc "scifi-search/app/database"
	"scifi-search/app/http/middlewares"
	"scifi-search/app/http/notifications/cookies"
	"scifi-search/app/languages"
	"scifi-search/app/utils"
	"scifi-search/app/utils/checkers"
	"scifi-search/app/views"

	"github.com/a-h/templ"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

// ---------------------------------------------------------------------
// Services
// ---------------------------------------------------------------------

func RegisterAuthenticationHandlers() {

	// Creación de la cuenta de administración.
	createAdmin()

	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/login", logInHandler)
	http.HandleFunc("/signout", signOutHandler)
	http.HandleFunc("/delete-account", deleteUserHandler)

	http.HandleFunc(
		"/loader",
		middlewares.RequiresEmailVerified(
			middlewares.RequiresAuthorization(
				loaderHandler,
				auth.AdminRole.Level,
			),
		),
	)

	// Handler de verificación de email.
	http.HandleFunc("/auth/verify-email", verifyEmailHandler)

	http.HandleFunc("/auth/session/refresh", func(w http.ResponseWriter, r *http.Request) {
		session.RefreshSession(r, w)
	})

	// TODO: refactorizar.
	// Al usar VerifySession, la sesión ya está garantizada y puesta en el contexto
	// si las cookies de sesión son válidas.
	http.HandleFunc("/auth/sessioninfo", func(w http.ResponseWriter, r *http.Request) {

		payload, err := auth.GetSessionInfo(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		component := views.SessionInfo(payload)
		templ.Handler(component).ServeHTTP(w, r)
	})
}

// ------------------------------------------------------------------------------------------------
// Handlers
// ------------------------------------------------------------------------------------------------

func signUpHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		component := views.SignUpPage("", auth.GetCurrentAuthorizationLevel(w, r), languages.GetTranslatorFromRequest(r))
		templ.Handler(component).ServeHTTP(w, r)
	case http.MethodPost:
		signUp(w, r)
	default:
		http.Error(w, MethodNotAllowedError.Error(), http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func logInHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		component := views.LoginPage(
			auth.GetCurrentAuthorizationLevel(w, r),
			languages.GetTranslatorFromRequest(r),
		)
		component.Render(r.Context(), w)

	case http.MethodPost:
		logIn(w, r)

	default:
		http.Error(w, MethodNotAllowedError.Error(), http.StatusMethodNotAllowed)
	}
}

// ------------------------------------------------------------------------------------------------

func signOutHandler(w http.ResponseWriter, r *http.Request) {

	auth.RevokeSession(w, r)

	cookies.AddFlashCookie(w, "Successful signout!")

	http.Redirect(w, r, "/", http.StatusFound)
}

// ------------------------------------------------------------------------------------------------

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, MethodNotAllowedError.Error(), http.StatusMethodNotAllowed)
		return
	}

	deleteUser(w, r)
}

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

func signUp(w http.ResponseWriter, r *http.Request) {

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
	if checkers.IsThereAnEmptyField(name, surname, email, password) {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}

	newUser, err := createUser(name, surname, email, password, auth.UserRole, languages.GetTranslatorFromRequest(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = auth.CreateSession(w, r, newUser.AuthID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookies.AddFlashCookie(w, languages.GetTranslatorFromRequest(r)("email-verification.sent"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------

func createAdmin() {

	// Obtención de los datos del administrador.
	name := utils.GetEnv("ADMIN_NAME", "admin")
	surname := utils.GetEnv("ADMIN_SURNAME", "full-access")
	email := utils.GetEnv("ADMIN_EMAIL", "admin@scifi-search.com")
	password := utils.GetEnv("ADMIN_PASSWORD", "admin")

	_, err := createUser(name, surname, email, password, auth.AdminRole, languages.GetTranslatorFromRequest(nil))
	if err != nil {
		if !errors.Is(err, auth.EmailAlreadyInUseError) {
			log.Fatal(UnknownError.Error())
		}
	}
}

// ------------------------------------------------------------------------------------------------

func createUser(name, surname, email, password string, role auth.Role, translator languages.Translator) (*sqlc.User, error) {

	userID, err := auth.RegisterUser(email, password)
	if err != nil {
		return nil, err
	}

	// Se crea el usuario en la base de datos.
	user, err := queries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		Name:    name,
		Surname: surname,
		AuthID:  *userID,
	})
	if err != nil {
		return nil, err
	}

	emailSubject := translator("verification-email.subject")
	emailBody := translator("verification-email.body")

	auth.SendVerificationEmail(emailService, *userID, email, emailSubject, emailBody)
	auth.AssignRoleToUser(role, *userID)

	return &user, nil
}

// ------------------------------------------------------------------------------------------------

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, "Token no proporcionado", http.StatusBadRequest)
		return
	}

	err := auth.VerifyEmail(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookies.AddFlashCookie(w, "¡Email verificado exitosamente!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------

func logIn(w http.ResponseWriter, r *http.Request) {
	translator := languages.GetTranslatorFromRequest(r)

	if r.Method == http.MethodGet {
		views.LoginPage(auth.GetCurrentAuthorizationLevel(w, r), translator).Render(r.Context(), w)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, MethodNotAllowedError.Error(), http.StatusMethodNotAllowed)
		return
	}

	// Se parsean y obtienen las credenciales.
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al parsear formulario", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Se validan las credenciales.
	userID, err := auth.VerifyCredentials(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Se intenta crear la sesión.
	err = auth.CreateSession(w, r, *userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookies.AddFlashCookie(w, "Welcome back!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------

func deleteUser(w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	userEmail := auth.GetCurrentUserEmail(w, r)
	if userEmail == nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	err = auth.DeleteUser(user.AuthID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Se elimina el avatar del usuario (si tiene).
	err = avatarService.Delete(r.Context(), user.UserID)
	if err != nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	// Se elimina el usuario de la base de datos.
	err = queries.DeleteUser(r.Context(), user.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: El usuario no existe.
			http.Error(w, UserNotFoundError.Error(), http.StatusNotFound)
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

	emailBody := strings.Join([]string{
		"¡Nos entristece ver que te vayas! \n",
		"Para tu seguridad y tranquilidad, hemos eliminado todos tus datos. \n",
		"¡Esperamos volver a verte pronto!",
	}, "")

	emailService.Send(
		*userEmail,
		"¡Lamentamos que te vayas!",
		emailBody,
	)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------

func getCurrentUser(w http.ResponseWriter, r *http.Request) (*database.User, error) {

	authID, err := auth.GetCurrentUserID(w, r)
	if err != nil {
		return nil, err
	}

	user, err := queries.GetUserByAuthID(r.Context(), *authID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ------------------------------------------------------------------------------------------------
