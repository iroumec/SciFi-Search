package auth

// --------------------------------------------------------------------------------------------- //
// Importaciones
// --------------------------------------------------------------------------------------------- //

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"scifi-search/app/database"
	"scifi-search/app/infra/auth/supertokens"
	"scifi-search/app/utils/converters"
	"scifi-search/app/workers"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
)

// --------------------------------------------------------------------------------------------- //
// Servicios Generales
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

func SendVerificationEmail(userID, email string) {

	tokenResponse, err := emailverification.CreateEmailVerificationToken("", userID, &email, nil)
	if err != nil {
		log.Println("Error creando token de verificación:", err)
	} else if tokenResponse.OK != nil {
		verificationLink := supertokens.WebsiteDomain + "/auth/verify-email?token=" + tokenResponse.OK.Token

		// Se construye el cuerpo del email.
		subject := "Verifica tu email"
		body := fmt.Sprintf("Por favor. Verifica tu email entrando en el siguiente enlace:\n%s", verificationLink)

		// Envío asíncrono del email (no bloquea la respuesta HTTP).
		workers.SendEmailAsync(email, subject, body)
	}
}

// --------------------------------------------------------------------------------------------- //

// Actualiza el email del usuario y solicita re-verificación si el email cambió.
func UpdateEmail(user *database.User, newEmail string) error {
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
	SendVerificationEmail(user.AuthID, newEmail)

	log.Println("Email actualizado exitosamente, se envió verificación")
	return nil
}

// --------------------------------------------------------------------------------------------- //

// Actualiza la contraseña del usuario.
func UpdatePassword(user *database.User, currentPassword, newPassword string) error {
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
