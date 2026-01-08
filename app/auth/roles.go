package auth

import (
	"scifi-search/app/infra/auth/supertokens"
)

// ---------------------------------------------------------------------
// Estructuras
// ---------------------------------------------------------------------

// Rol del sistema, con su nombre y nivel de autorización.
type Role struct {
	Name  string
	Level int
}

// ---------------------------------------------------------------------
//  Variables
// ---------------------------------------------------------------------

// Roles predefinidos del sistema.
var (
	NoRole     = Role{Name: "no-role", Level: -1}
	AdminRole  = Role{Name: "admin", Level: 2}
	LoaderRole = Role{Name: "loader", Level: 1}
	UserRole   = Role{Name: "user", Level: 0}
)

// ---------------------------------------------------------------------
//  Funciones
// ---------------------------------------------------------------------

func AssignRoleToUser(role Role, userID string) {

	supertokens.AssignRoleToUser(userID, role.Name)
}

// ---------------------------------------------------------------------
