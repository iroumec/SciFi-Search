package handlers

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

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

// ------------------------------------------------------------------------------------------------
// Variables
// ------------------------------------------------------------------------------------------------

// Errors.
var (
	TokenNotProvidedError = errors.New("errors.token-not-provided")
)

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

func RegisterAuthenticationHandlers() {

	createAdministratorAccount()

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

	http.HandleFunc("/auth/verify-email", verifyEmailHandler)

	http.HandleFunc("/auth/session/refresh", func(w http.ResponseWriter, r *http.Request) {
		session.RefreshSession(r, w)
	})

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
		component := views.SignUpPage(
			"",
			auth.GetCurrentAuthorizationLevel(w, r),
			languages.GetTranslatorFromRequest(r),
		)
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

	cookies.AddFlashCookie(w, "messages.sign-out")

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
// Functions
// ------------------------------------------------------------------------------------------------

func signUp(w http.ResponseWriter, r *http.Request) {

	// Form parsing.
	if err := r.ParseForm(); err != nil {
		http.Error(w, FormParsingError.Error(), http.StatusBadRequest)
		return
	}

	// User data.
	name := r.Form.Get("name")
	surname := r.Form.Get("surname")
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Validation.
	if checkers.IsThereAnEmptyField(name, surname, email, password) {
		http.Error(w, MissingRequiredFieldsError.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := createUser(
		name,
		surname,
		email,
		password,
		auth.UserRole,
		languages.GetTranslatorFromRequest(r),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = auth.CreateSession(w, r, newUser.AuthID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookies.AddFlashCookie(w, languages.GetTranslatorFromRequest(r)("verification-email-sent"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------

func createAdministratorAccount() {

	// Administrator data.
	name := utils.GetEnv("ADMIN_NAME", "admin")
	surname := utils.GetEnv("ADMIN_SURNAME", "full-access")
	email := utils.GetEnv("ADMIN_EMAIL", "admin@scifi-search.com")
	password := utils.GetEnv("ADMIN_PASSWORD", "admin")

	_, err := createUser(
		name,
		surname,
		email,
		password,
		auth.AdminRole,
		languages.GetTranslatorFromRequest(nil),
	)
	if err != nil {
		if !errors.Is(err, auth.EmailAlreadyInUseError) {
			log.Fatal(UnknownError.Error())
		}
	}
}

// ------------------------------------------------------------------------------------------------

func createUser(
	name, surname, email, password string, role auth.Role, translator languages.Translator,
) (*sqlc.User, error) {

	userID, err := auth.RegisterUser(email, password)
	if err != nil {
		return nil, err
	}

	// User is created in the database.
	user, err := queries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		Name:    name,
		Surname: surname,
		AuthID:  *userID,
	})
	if err != nil {
		return nil, err
	}

	emailSubject := translator("messages.verification-email-subject")
	emailBody := translator("messages.verification-email-body")

	auth.SendVerificationEmail(emailService, *userID, email, emailSubject, emailBody)
	auth.AssignRoleToUser(role, *userID)

	return &user, nil
}

// ------------------------------------------------------------------------------------------------

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(w, TokenNotProvidedError.Error(), http.StatusBadRequest)
		return
	}

	err := auth.VerifyEmail(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookies.AddFlashCookie(w, languages.GetTranslatorFromRequest(r)("messages.email-verified"))
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

	// Form parsing.
	if err := r.ParseForm(); err != nil {
		http.Error(w, FormParsingError.Error(), http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Credentials validation.
	userID, err := auth.VerifyCredentials(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Session creation.
	err = auth.CreateSession(w, r, *userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookies.AddFlashCookie(w, "messages.welcome-back")
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

	// User's avatar is deleted.
	err = avatarService.Delete(r.Context(), user.UserID)
	if err != nil {
		http.Error(w, InternalServerError.Error(), http.StatusInternalServerError)
		return
	}

	// User is deleted from the database.
	err = queries.DeleteUser(r.Context(), user.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Error 404: User doesn't exist.
			http.Error(w, UserNotFoundError.Error(), http.StatusNotFound)
		} else {
			// Error 500: Unexpected Internal Error.
			http.Error(w, UnexpectedError.Error(), http.StatusInternalServerError)
		}
		return
	}

	auth.RevokeSession(w, r)

	translator := languages.GetTranslatorFromRequest(r)

	cookies.AddFlashCookie(w, translator("messages.user-deleted"))

	emailBody := strings.Join([]string{
		translator("messages.user-deletion-email-body-part-1") + "\n\n",
		translator("messages.user-deletion-email-body-part-2") + "\n\n",
		translator("messages.user-deletion-email-body-part-3") + "\n\n",
	}, "")

	emailService.Send(
		*userEmail,
		translator("messages.user-deletion-email-subject"),
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
