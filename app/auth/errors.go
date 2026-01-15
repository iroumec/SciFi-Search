package auth

// ------------------------------------------------------------------------------------------------
//  Imports
// ------------------------------------------------------------------------------------------------

import "errors"

// ------------------------------------------------------------------------------------------------
//  Variables
// ------------------------------------------------------------------------------------------------

var (

	// Authentication of user
	ErrNotAuthenticated         = errors.New("error.user-not-authenticated")
	UserNotFoundError           = errors.New("error.user-not-found")
	WrongCredentialsError       = errors.New("error.wrong-credentials")
	EmailAlreadyInUseError      = errors.New("error.email-already-in-use")
	RequiredDataNotSpecified    = errors.New("error.required-data-not-specified")
	PasswordPolicyViolatedError = errors.New("error.password-policy-violated")
	Unauthorized                = errors.New("error.unauthorized") // 401 -> No authenticated
	Forbidden                   = errors.New("error.forbidden")    // 403 -> No authorized

	// Validation of email
	InvalidOrExpiredTokenError = errors.New("error.invalid-or-expired-token")

	// Session
	NoSessionError = errors.New("error.no-session")

	// Generals
	UnknownError = errors.New("error.unknown")
)

// ------------------------------------------------------------------------------------------------
