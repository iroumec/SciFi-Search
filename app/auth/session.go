package auth

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"errors"
	"fmt"
	"net/http"
	"scifi-search/app/utils/converters"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

// ------------------------------------------------------------------------------------------------
// Servicios
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

func RevokeSession(w http.ResponseWriter, r *http.Request) error {

	sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
		SessionRequired: converters.ToBoolPointer(false), // False -> No error si no hay sessión. Posiblemente deba cambiarse luego.
	})
	if err != nil {
		return Unauthorized
	}

	if sessionContainer == nil {
		return NoSessionError
	}

	if err := sessionContainer.RevokeSession(); err != nil {
		return err
	}

	return nil
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

	session, err := getCurrentSession(w, r)
	if err != nil {
		return nil, err
	}

	userID := session.GetUserID()

	return &userID, nil
}

// ------------------------------------------------------------------------------------------------

func GetSessionInfo(w http.ResponseWriter, r *http.Request) (map[string]string, error) {

	session, err := getCurrentSession(w, r)
	if err != nil {
		return nil, err
	}

	rawPayload := session.GetAccessTokenPayload()

	payload := make(map[string]string)
	for k, v := range rawPayload {
		payload[k] = fmt.Sprintf("%v", v)
	}

	return payload, nil
}

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

func getCurrentSession(w http.ResponseWriter, r *http.Request) (sessmodels.SessionContainer, error) {

	sessionContainer, err := session.GetSession(r, w, &sessmodels.VerifySessionOptions{
		SessionRequired: converters.ToBoolPointer(false),
	})
	if err != nil {
		return nil, NoSessionError
	}

	if sessionContainer == nil {
		return nil, NoSessionError
	}

	return sessionContainer, nil
}

// ------------------------------------------------------------------------------------------------
