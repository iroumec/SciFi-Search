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
	ErrNotAuthenticated         = errors.New("errors.user-not-authenticated")
	UserNotFoundError           = errors.New("errors.user-not-found")
	WrongCredentialsError       = errors.New("errors.wrong-credentials")
	EmailAlreadyInUseError      = errors.New("errors.email-already-in-use")
	RequiredDataNotSpecified    = errors.New("errors.required-data-not-specified")
	PasswordPolicyViolatedError = errors.New("errors.password-policy-violated")
	Unauthorized                = errors.New("errors.unauthorized") // 401 -> Not authenticated
	Forbidden                   = errors.New("errors.forbidden")    // 403 -> Not authorized

	// Validation of email
	InvalidOrExpiredTokenError = errors.New("errors.invalid-or-expired-token")

	// Session
	NoSessionError = errors.New("errors.no-session")

	// Generals
	UnknownError = errors.New("errors.unknown")
)

// ------------------------------------------------------------------------------------------------
