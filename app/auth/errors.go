package auth

// ------------------------------------------------------------------------------------------------
//  Importaciones
// ------------------------------------------------------------------------------------------------

import "errors"

// ------------------------------------------------------------------------------------------------
//  Variables
// ------------------------------------------------------------------------------------------------

var (

	// Authentication of user
	ErrNotAuthenticated         = errors.New("user not authenticated")
	UserNotFoundError           = errors.New("User not found")
	WrongCredentialsError       = errors.New("Wrong credentials")
	EmailAlreadyInUseError      = errors.New("Email already in use")
	RequiredDataNotSpecified    = errors.New("Required data was not specified")
	PasswordPolicyViolatedError = errors.New("Password policy violated")
	Unauthorized                = errors.New("Unauthorized") // 401 -> No autenticado
	Forbidden                   = errors.New("Forbidden")    // 403 -> No autorizado

	// Validation of email
	ErrInvalidOrExpiredToken = errors.New("invalid or expired token")

	// Generales
	UnknownError = errors.New("Unknown error")

	NoSessionError = errors.New("No session")
)

// ------------------------------------------------------------------------------------------------
