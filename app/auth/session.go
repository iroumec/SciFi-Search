package auth

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"scifi-search/app/database"
	"scifi-search/app/utils/converters"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

// ------------------------------------------------------------------------------------------------
// Servicios (ligados a la sesión actual)
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

func GetCurrentAuthorizationLevel(w http.ResponseWriter, r *http.Request) int {

	authorizationLevel := NoRole.Level
	authID, err := GetCurrentUserID(w, r)
	if err != nil {
		if !errors.Is(err, ErrNotAuthenticated) {
			return NoRole.Level
		}
	} else {
		authorizationLevel = GetAuthenticationLevel(*authID)
	}

	return authorizationLevel
}

// ------------------------------------------------------------------------------------------------

func GetCurrentUserEmail(w http.ResponseWriter, r *http.Request) *string {

	if IsUserAuthenticated(w, r) {

		sessionContainer, _ := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
			SessionRequired: converters.ToBoolPointer(false),
		})

		return GetUserEmail(sessionContainer.GetUserID())

	} else {

		return nil
	}
}

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

func CreateSession(w http.ResponseWriter, r *http.Request, userID string) error {

	_, err := session.CreateNewSession(r, w, "", userID, nil, nil)
	if err != nil {
		return UnknownError
	}

	return nil
}

// ------------------------------------------------------------------------------------------------

func GetCurrentUserID(w http.ResponseWriter, r *http.Request) (*string, error) {

	sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
		SessionRequired: converters.ToBoolPointer(false),
	})
	if err != nil {
		return nil, err
	}

	if sessionContainer == nil {
		return nil, NoSessionError
	}

	userID := sessionContainer.GetUserID()

	return &userID, nil
}

// ------------------------------------------------------------------------------------------------
