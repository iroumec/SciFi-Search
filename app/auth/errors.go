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
var ErrUserNotFound = errors.New("user not found")
var WrongCredentialsError = errors.New("Wrong credentials")
var EmailAlreadyInUseError = errors.New("Email already in use")

// Validation of email
var ErrInvalidOrExpiredToken = errors.New("invalid or expired token")

// Generales
var UnknownError = errors.New("Unknown error")

// ---------------------------------------------------------------------
