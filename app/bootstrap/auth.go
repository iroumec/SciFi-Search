package bootstrap

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import (
	"scifi-search/app/auth"
	"scifi-search/app/email"
	"scifi-search/app/infra/auth/supertokens"
	"scifi-search/app/infra/email/mailhog"
)

// ------------------------------------------------------------------------------------------------
// Functions
// ------------------------------------------------------------------------------------------------

func initializeAuthorizationMechanisms() {

	emailProvider := mailhog.New()
	emailService := email.New(emailProvider)

	supertokens.Initialize(emailService)

	defineRolesInProvider()
}

// ------------------------------------------------------------------------------------------------

// Roles creations.
func defineRolesInProvider() {

	supertokens.CreateNewRoleOrAddPermissions(
		auth.AdminRole.Name,
		[]string{"full-access"},
	) // Administrator.

	supertokens.CreateNewRoleOrAddPermissions(
		auth.LoaderRole.Name,
		[]string{"manage-own-financings"},
	) // Loaders.

	supertokens.CreateNewRoleOrAddPermissions(
		auth.UserRole.Name,
		[]string{"view-only"},
	) // Users.
}

// ------------------------------------------------------------------------------------------------
