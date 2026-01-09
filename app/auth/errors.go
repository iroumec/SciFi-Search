package auth

// ---------------------------------------------------------------------
//  Importaciones
// ---------------------------------------------------------------------

import "errors"

// ---------------------------------------------------------------------
//  Variables
// ---------------------------------------------------------------------

// Authentication of user
var ErrNotAuthenticated = errors.New("user not authenticated")
var UserNotFoundError = errors.New("User not found")
var WrongCredentialsError = errors.New("Wrong credentials")
var EmailAlreadyInUseError = errors.New("Email already in use")
var RequiredDataNotSpecified = errors.New("Required data was not specified")
var PasswordPolicyViolatedError = errors.New("Password policy violated")
var Unauthorized = errors.New("Unauthorized") // 401 -> No autenticado
var Forbidden = errors.New("Forbidden")       // 403 -> No autorizado

// Validation of email
var ErrInvalidOrExpiredToken = errors.New("invalid or expired token")

// Generales
var UnknownError = errors.New("Unknown error")

var NoSessionError = errors.New("No session")

// ---------------------------------------------------------------------
