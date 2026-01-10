package bootstrap

// ------------------------------------------------------------------------------------------------
// Importaciones
// ------------------------------------------------------------------------------------------------

import (
	"scifi-search/app/auth"
	"scifi-search/app/email"
	"scifi-search/app/infra/auth/supertokens"
	"scifi-search/app/infra/email/mailhog"
)

// ------------------------------------------------------------------------------------------------
// Funciones
// ------------------------------------------------------------------------------------------------

func initializeAuthorizationMechanisms() {

	mailProvider := mailhog.New()
	emailService := email.New(mailProvider)

	supertokens.Initialize(emailService)

	defineRolesInProvider()
}

// ------------------------------------------------------------------------------------------------

func defineRolesInProvider() {

	// Creación de roles.
	supertokens.CreateNewRoleOrAddPermissions(
		auth.AdminRole.Name,
		[]string{"full-access"},
	) // Administrador.

	supertokens.CreateNewRoleOrAddPermissions(
		auth.LoaderRole.Name,
		[]string{"manage-own-financings"},
	) // Entidades que cargan financiamiento.

	supertokens.CreateNewRoleOrAddPermissions(
		auth.UserRole.Name,
		[]string{"view-only"},
	) // Usuarios normales.
}

// ------------------------------------------------------------------------------------------------
