package auth

// --------------------------------------------------------------------------------------------- //
// Imports
// --------------------------------------------------------------------------------------------- //

import (
	"fmt"
	"slices"
	"strings"

	"scifi-search/app/email"
	"scifi-search/app/infra/auth/supertokens"
	"scifi-search/app/utils/converters"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
)

// --------------------------------------------------------------------------------------------- //
// Constants
// --------------------------------------------------------------------------------------------- //

const (
	tokenVerificationPath = "/auth/verify-email?token="
)

// --------------------------------------------------------------------------------------------- //
// Services
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

	// Unauthenticated user.
	return NoRole.Level
}

// --------------------------------------------------------------------------------------------- //

func SendVerificationEmail(emailService *email.Service, userID, email, emailSubject, emailBody string) error {

	tokenResponse, err := emailverification.CreateEmailVerificationToken("", userID, &email, nil)
	if err != nil {
		return err
	} else if tokenResponse.OK != nil {

		verificationLink := supertokens.WebsiteDomain + tokenVerificationPath + tokenResponse.OK.Token
		body := fmt.Sprintf("%s\n%s", emailBody, verificationLink)

		emailService.Send(email, emailSubject, body)
	}

	return nil
}

// --------------------------------------------------------------------------------------------- //

// Updates the user's email and asks for re-verification if the email was changed.
func UpdateEmail(authID, newEmail string) error {
	newEmail = strings.TrimSpace(newEmail)

	// Obtaining the current user's email
	currentUser, err := emailpassword.GetUserByID(authID)
	if err != nil || currentUser == nil {
		return err
	}

	currentEmail := currentUser.Email

	// If there's no change, nothing is made.
	if strings.EqualFold(currentEmail, newEmail) {
		return nil
	}

	// If there exists a user with that email.
	existingUser, err := emailpassword.GetUserByEmail("", newEmail)
	if err != nil {
		return err
	}

	if existingUser != nil && existingUser.ID != authID {
		return EmailAlreadyInUseError
	}

	// Email update on SuperTokens.
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

	// Email unverification (forcing re-verification).
	_, err = emailverification.UnverifyEmail(authID, &newEmail, nil)
	if err != nil {
		return err
	}

	return nil
}

// --------------------------------------------------------------------------------------------- //

// Updates the user's password.
func UpdatePassword(authID, currentPassword, newPassword string) error {
	// Checking if both password were given.
	if currentPassword == "" || newPassword == "" {
		return nil
	}

	// Obtaining email to verify the password.
	currentUser, err := emailpassword.GetUserByID(authID)
	if err != nil || currentUser == nil {
		return err
	}

	_, err = VerifyCredentials(currentUser.Email, currentPassword)
	if err != nil {
		return err
	}

	// Password update.
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

	// User deletion in Supertokens
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

	// Email already in use validation.
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

	// Token verification.
	response, err := emailverification.VerifyEmailUsingToken("", token, nil)
	if err != nil {
		return err
	}

	if response.OK != nil {
		return nil
	} else if response.EmailVerificationInvalidTokenError != nil {
		return InvalidOrExpiredTokenError
	} else {
		return UnknownError
	}
}

// --------------------------------------------------------------------------------------------- //

// Returns the user's ID matching the credentials.
func VerifyCredentials(email, password string) (*string, error) {

	// Log-in attempt.
	resp, err := emailpassword.SignIn("", email, password)
	if err != nil {
		return nil, err
	}

	if resp.WrongCredentialsError != nil {
		return nil, WrongCredentialsError
	}

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
