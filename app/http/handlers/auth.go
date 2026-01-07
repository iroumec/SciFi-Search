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
	"strings"

	"scifi-search/app/database"
	sqlc "scifi-search/app/database"
	"scifi-search/app/http/cookies"
	"scifi-search/app/infra/auth"
	"scifi-search/app/languages"
	"scifi-search/app/utils"
	"scifi-search/app/utils/checkers"
	"scifi-search/app/utils/converters"
	"scifi-search/app/views"
	"scifi-search/app/workers"

	"github.com/a-h/templ"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// ---------------------------------------------------------------------
// Constantes
// ---------------------------------------------------------------------

const (
	websiteDomain = "http://localhost:8080"
)

// ---------------------------------------------------------------------

var ErrNotAuthenticated = errors.New("user not authenticated")
var ErrUserNotFound = errors.New("user not found")

// ---------------------------------------------------------------------

func registerAuthenticationHandlers() {

	auth.InitializeSupertokens()

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
		component := views.SignUpPage("", languages.GetTranslatorFromRequest(r))
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

func loaderHandler(w http.ResponseWriter, r *http.Request) {

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

}

// ------------------------------------------------------------------------------------------------

func createLoader(name, surname, email, password string) {

	user, resp := createUser(name, surname, email, password, auth.LoaderRole.Name)

	if user == nil || resp == nil {
		log.Fatal("Ocurrió un error al momento de crear al usuario.")
	}

	log.Printf("Loader created: %s", email)
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
		sendVerificationEmail(resp.OK.User.ID, email)
	}

	// Asignación del rol al usuario.
	userroles.AddRoleToUser("public", resp.OK.User.ID, role, nil)

	return &newUser, &resp
}

// ------------------------------------------------------------------------------------------------

func sendVerificationEmail(userID, email string) {

	tokenResponse, err := emailverification.CreateEmailVerificationToken("", userID, &email, nil)
	if err != nil {
		log.Println("Error creando token de verificación:", err)
	} else if tokenResponse.OK != nil {
		verificationLink := websiteDomain + "/auth/verify-email?token=" + tokenResponse.OK.Token

		// Se construye el cuerpo del email.
		subject := "Verifica tu email"
		body := fmt.Sprintf("Por favor. Verifica tu email entrando en el siguiente enlace:\n%s", verificationLink)

		// Envío asíncrono del email (no bloquea la respuesta HTTP).
		workers.SendEmailAsync(email, subject, body)
	}
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
		views.LoginPage(translator).Render(r.Context(), w)
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
		views.LoginPage(translator).Render(r.Context(), w)
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

	revokeSession(w, r)

	cookies.AddFlashCookie(w, "Successful signout!")

	http.Redirect(w, r, "/", http.StatusFound)
}

// ------------------------------------------------------------------------------------------------
// Funciones Auxiliares
// ------------------------------------------------------------------------------------------------

// Retorna si el usuario está autenticado.
func isUserAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	// Intentar obtener la sesión sin requerirla
	sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
		SessionRequired: converters.ToBoolPointer(false),
	})

	return err == nil && sessionContainer != nil
}

// ------------------------------------------------------------------------------------------------
// Get Current User
// ------------------------------------------------------------------------------------------------

func getCurrentUser(w http.ResponseWriter, r *http.Request) (*database.User, error) {
	if !isUserAuthenticated(w, r) {
		return nil, ErrNotAuthenticated
	}

	sessionContainer, err := session.GetSession(r, nil, &sessmodels.VerifySessionOptions{
		SessionRequired: converters.ToBoolPointer(false),
	})
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	supertokensUserID := sessionContainer.GetUserID()

	user, err := queries.GetUserByAuthID(r.Context(), supertokensUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by auth id: %w", err)
	}

	return &user, nil
}

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

	user, err := getCurrentUser(w, r)
	if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	userEmail := getCurrentUserEmail(w, r)
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
	err = deleteAvatar(r.Context(), user.UserID)
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
	revokeSession(w, r)

	cookies.AddFlashCookie(w, "Usuario eliminado. ¡Lamentamos que te vayas!")

	workers.SendEmailAsync(
		*userEmail,
		"¡Lamentamos que te vayas!",
		"Nos entristece ver que te vayas. Para tu seguridad, hemos eliminado todos tus datos. ¡Esperamos volver a verte pronto!",
	)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ------------------------------------------------------------------------------------------------

func revokeSession(w http.ResponseWriter, r *http.Request) {

	sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
		SessionRequired: converters.ToBoolPointer(false), // False -> No error si no hay sessión. Posiblemente deba cambiarse luego.
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

func getUserEmail(userID string) *string {

	user, err := emailpassword.GetUserByID(userID)
	if err != nil || user == nil {
		return nil
	}

	return &user.Email
}

// ------------------------------------------------------------------------------------------------

func getCurrentUserEmail(w http.ResponseWriter, r *http.Request) *string {

	if isUserAuthenticated(w, r) {

		sessionContainer, _ := session.GetSession(r, nil, &sessmodels.VerifySessionOptions{
			SessionRequired: converters.ToBoolPointer(false),
		})

		return getUserEmail(sessionContainer.GetUserID())

	} else {

		return nil
	}
}

// ------------------------------------------------------------------------------------------------

// Actualiza el email del usuario y solicita re-verificación si el email cambió.
func updateEmail(user *database.User, newEmail string) error {
	newEmail = strings.TrimSpace(newEmail)

	// Obtención del email actual del usuario.
	currentUser, err := emailpassword.GetUserByID(user.AuthID)
	if err != nil || currentUser == nil {
		log.Println("Error al obtener usuario actual:", err)
		return err
	}

	currentEmail := currentUser.Email

	// Si el email es el mismo, no se hace nada.
	if strings.EqualFold(currentEmail, newEmail) {
		log.Println("El email no cambió, omitiendo actualización")
		return nil
	}

	// Se verifica si el nuevo emaul ya está en uso por otro usuario.
	existingUser, err := emailpassword.GetUserByEmail("", newEmail)
	if err != nil {
		log.Println("Error al verificar email existente:", err)
		return err
	}

	if existingUser != nil && existingUser.ID != user.AuthID {
		log.Println("El email ya está en uso por otro usuario")
		return err
	}

	// Se actualiza el email en SuperTokens.
	updateResp, err := emailpassword.UpdateEmailOrPassword(user.AuthID, &newEmail, nil, nil, nil)
	if err != nil {
		log.Println("Error actualizando email:", err)
		return err
	}

	if updateResp.OK == nil {
		if updateResp.EmailAlreadyExistsError != nil {
			log.Println("Email ya existe")
		} else if updateResp.UnknownUserIdError != nil {
			log.Println("Usuario no encontrado")
		}
		return err
	}

	// Se desverifica el email (forzando re-verificación).
	_, err = emailverification.UnverifyEmail(user.AuthID, &newEmail, nil)
	if err != nil {
		log.Println("Error al desverificar email:", err)
	}

	// Se envía un email de verificación con un nuevo token al nuevo email.
	sendVerificationEmail(user.AuthID, newEmail)

	log.Println("Email actualizado exitosamente, se envió verificación")
	return nil
}

// ------------------------------------------------------------------------------------------------

// Actualiza la contraseña del usuario.
func updatePassword(user *database.User, currentPassword, newPassword string) error {
	// Se valida que se proporcionaron ambas contraseñas.
	if currentPassword == "" || newPassword == "" {
		log.Println("Contraseñas no proporcionadas")
		return nil
	}

	// Se obtiene el email actual para verificar la contraseña.
	currentUser, err := emailpassword.GetUserByID(user.AuthID)
	if err != nil || currentUser == nil {
		log.Println("Error al obtener usuario actual:", err)
		return err
	}

	// Se verifica la contraseña actual.
	signInResp, err := emailpassword.SignIn("", currentUser.Email, currentPassword)
	if err != nil {
		log.Println("Error al verificar contraseña:", err)
		return err
	}

	if signInResp.WrongCredentialsError != nil {
		log.Println("Contraseña actual incorrecta")
		return err
	}

	// Se actualiza la contraseña.
	updateResp, err := emailpassword.UpdateEmailOrPassword(user.AuthID, nil, &newPassword, nil, nil)
	if err != nil {
		log.Println("Error actualizando contraseña:", err)
		return err
	}

	if debug {

		// Ver qué tiene la respuesta.
		log.Printf("UpdateEmailOrPassword response - OK: %v, UnknownUser: %v, EmailExists: %v, PasswordPolicy: %v",
			updateResp.OK != nil,
			updateResp.UnknownUserIdError != nil,
			updateResp.EmailAlreadyExistsError != nil,
			updateResp.PasswordPolicyViolatedError != nil)

		if updateResp.PasswordPolicyViolatedError != nil {
			log.Printf("DETALLE del error de política: %+v", updateResp.PasswordPolicyViolatedError)
		}

	}

	if updateResp.OK == nil {
		if updateResp.UnknownUserIdError != nil {
			log.Println("Usuario no encontrado")
			return fmt.Errorf("usuario no encontrado")
		} else if updateResp.EmailAlreadyExistsError != nil {
			log.Println("Email ya existe (esto no debería ocurrir al cambiar solo contraseña)")
			return fmt.Errorf("email ya existe")
		} else if updateResp.PasswordPolicyViolatedError != nil {
			log.Println("La contraseña no cumple con la política de seguridad")
			return fmt.Errorf("contraseña no cumple con los requisitos de seguridad")
		} else {
			log.Println("Error desconocido al actualizar contraseña")
			return fmt.Errorf("error desconocido al actualizar contraseña")
		}
	}

	log.Println("Contraseña actualizada")
	return nil
}

// ------------------------------------------------------------------------------------------------

func getCurrentAuthorizationLevel(w http.ResponseWriter, r *http.Request) int {

	authID := auth.NoRole.Level
	currentUser, err := getCurrentUser(w, r)
	if err != nil {
		if !errors.Is(err, ErrNotAuthenticated) {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return auth.NoRole.Level
		}
	} else {
		authID = auth.GetAuthenticationLevel(currentUser.AuthID)
	}

	return authID
}

// ------------------------------------------------------------------------------------------------
