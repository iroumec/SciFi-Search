package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"scifi-search/app/database"
	"scifi-search/app/utils/converters"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

// ------------------------------------------------------------------------------------------------

var ErrNotAuthenticated = errors.New("user not authenticated")
var ErrUserNotFound = errors.New("user not found")

// ------------------------------------------------------------------------------------------------

func GetCurrentUser(w http.ResponseWriter, r *http.Request, queries *database.Queries) (*database.User, error) {
	if !IsUserAuthenticated(w, r) {
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

func GetCurrentAuthorizationLevel(w http.ResponseWriter, r *http.Request, queries *database.Queries) int {

	authID := NoRole.Level
	currentUser, err := GetCurrentUser(w, r, queries)
	if err != nil {
		if !errors.Is(err, ErrNotAuthenticated) {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return NoRole.Level
		}
	} else {
		authID = GetAuthenticationLevel(currentUser.AuthID)
	}

	return authID
}

// ------------------------------------------------------------------------------------------------

func GetCurrentUserEmail(w http.ResponseWriter, r *http.Request) *string {

	if IsUserAuthenticated(w, r) {

		sessionContainer, _ := session.GetSession(r, nil, &sessmodels.VerifySessionOptions{
			SessionRequired: converters.ToBoolPointer(false),
		})

		return GetUserEmail(sessionContainer.GetUserID())

	} else {

		return nil
	}
}

// ------------------------------------------------------------------------------------------------

// ------------------------------------------------------------------------------------------------
// Funciones Auxiliares
// ------------------------------------------------------------------------------------------------

// Retorna si el usuario está autenticado.
func IsUserAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	// Intentar obtener la sesión sin requerirla
	sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
		SessionRequired: converters.ToBoolPointer(false),
	})

	return err == nil && sessionContainer != nil
}

// ------------------------------------------------------------------------------------------------

func RevokeSession(w http.ResponseWriter, r *http.Request) {

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

func GetUserEmail(userID string) *string {

	user, err := emailpassword.GetUserByID(userID)
	if err != nil || user == nil {
		return nil
	}

	return &user.Email
}
