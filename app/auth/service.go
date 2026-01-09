package auth

// --------------------------------------------------------------------------------------------- //
// Importaciones
// --------------------------------------------------------------------------------------------- //

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"scifi-search/app/email"
	"scifi-search/app/infra/auth/supertokens"
	"scifi-search/app/utils/converters"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
)

// --------------------------------------------------------------------------------------------- //
// Servicios
// --------------------------------------------------------------------------------------------- //

func GetAuthenticationLevel(userID string) int {

	roles := supertokens.GetRolesForUser(userID)

	if slices.Contains(roles, AdminRole.Name) {
		return AdminRole.Level
	} else if slices.Contains(roles, LoaderRole.Name) {
		return LoaderRole.Level
	} else if slices.Contains(roles, UserRole.Name) {
		return UserRole.Level
	}

	// Usuario sin autenticar.
	return NoRole.Level
}

// --------------------------------------------------------------------------------------------- //

func SendVerificationEmail(emailService *email.Service, userID, email string) {

	tokenResponse, err := emailverification.CreateEmailVerificationToken("", userID, &email, nil)
	if err != nil {
		log.Println("Error creando token de verificación:", err)
		// TODO: implementar reintento.
	} else if tokenResponse.OK != nil {
		verificationLink := supertokens.WebsiteDomain + "/auth/verify-email?token=" + tokenResponse.OK.Token

		// Se construye el cuerpo del email.
		subject := "Verifica tu email"
		body := fmt.Sprintf("Por favor. Verifica tu email entrando en el siguiente enlace:\n%s", verificationLink)

		// Envío asíncrono del email (no bloquea la respuesta HTTP).
		emailService.Send(email, subject, body)
	}
}

// --------------------------------------------------------------------------------------------- //

// Actualiza el email del usuario y solicita re-verificación si el email cambió.
func UpdateEmail(authID, newEmail string) error {
	newEmail = strings.TrimSpace(newEmail)

	// Obtención del email actual del usuario.
	currentUser, err := emailpassword.GetUserByID(authID)
	if err != nil || currentUser == nil {
		return err
	}

	currentEmail := currentUser.Email

	// Si el email es el mismo, no se hace nada.
	if strings.EqualFold(currentEmail, newEmail) {
		return nil
	}

	// Se verifica si el nuevo emaul ya está en uso por otro usuario.
	existingUser, err := emailpassword.GetUserByEmail("", newEmail)
	if err != nil {
		return err
	}

	if existingUser != nil && existingUser.ID != authID {
		return EmailAlreadyInUseError
	}

	// Se actualiza el email en SuperTokens.
	updateResp, err := emailpassword.UpdateEmailOrPassword(authID, &newEmail, nil, nil, nil)
	if err != nil {
		return err
	}

	if updateResp.OK == nil {
		if updateResp.EmailAlreadyExistsError != nil {
			return EmailAlreadyInUseError
		} else if updateResp.UnknownUserIdError != nil {
			return UserNotFoundError
		}
		return err
	}

	// Se desverifica el email (forzando re-verificación).
	_, err = emailverification.UnverifyEmail(authID, &newEmail, nil)
	if err != nil {
		return err
	}

	return nil
}

// --------------------------------------------------------------------------------------------- //

// Actualiza la contraseña del usuario.
func UpdatePassword(authID, currentPassword, newPassword string) error {
	// Se valida que se proporcionaron ambas contraseñas.
	if currentPassword == "" || newPassword == "" {
		return nil
	}

	// Se obtiene el email actual para verificar la contraseña.
	currentUser, err := emailpassword.GetUserByID(authID)
	if err != nil || currentUser == nil {
		return err
	}

	_, err = VerifyCredentials(currentUser.Email, currentPassword)
	if err != nil {
		return err
	}

	// Se actualiza la contraseña.
	updateResp, err := emailpassword.UpdateEmailOrPassword(authID, nil, &newPassword, nil, nil)
	if err != nil {
		return err
	}

	if updateResp.OK == nil {
		if updateResp.UnknownUserIdError != nil {
			return UserNotFoundError
		} else if updateResp.EmailAlreadyExistsError != nil {
			return EmailAlreadyInUseError
		} else if updateResp.PasswordPolicyViolatedError != nil {
			return PasswordPolicyViolatedError
		} else {
			return UnknownError
		}
	}

	return nil
}

// --------------------------------------------------------------------------------------------- //

func GetUserEmail(userID string) *string {

	user, err := emailpassword.GetUserByID(userID)
	if err != nil || user == nil {
		return nil
	}

	return &user.Email
}

// --------------------------------------------------------------------------------------------- //

func DeleteUser(userID string) error {

	// Se elimina el usuario de supertokens.
	err := supertokens.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}

// --------------------------------------------------------------------------------------------- //

func RegisterUser(email, password string) (*string, error) {

	resp, err := emailpassword.SignUp("", email, password)
	if err != nil {
		return nil, UnknownError
	}

	// Verificación de email ya registrado.
	if resp.EmailAlreadyExistsError != nil {
		return nil, EmailAlreadyInUseError
	}

	if resp.OK != nil {
		return &resp.OK.User.ID, nil
	}

	return nil, UnknownError
}

// --------------------------------------------------------------------------------------------- //

func VerifyEmail(token string) error {

	// Se verifica el token.
	response, err := emailverification.VerifyEmailUsingToken("", token, nil)
	if err != nil {
		return err
	}

	if response.OK != nil {
		return nil
	} else if response.EmailVerificationInvalidTokenError != nil {
		return ErrInvalidOrExpiredToken
	} else {
		return UnknownError
	}
}

// --------------------------------------------------------------------------------------------- //

// Retorna la ID del usuario al que corresponden las credenciales.
func VerifyCredentials(email, password string) (*string, error) {

	// Se intenta realizar un log in.
	resp, err := emailpassword.SignIn("", email, password)
	if err != nil {
		return nil, err
	}

	// Credenciales incorrectas.
	if resp.WrongCredentialsError != nil {
		return nil, WrongCredentialsError
	}

	// Login exitoso: se crea la sesión y se redirige.
	if resp.OK != nil {
		return &resp.OK.User.ID, nil
	}

	return nil, UnknownError
}

// --------------------------------------------------------------------------------------------- //

func IsEmailVerified(authID string) (*bool, error) {

	isVerified, err := emailverification.IsEmailVerified(authID, nil, nil)
	if err != nil {
		return nil, err
	}

	return converters.ToBoolPointer(isVerified), nil
}

// --------------------------------------------------------------------------------------------- //
