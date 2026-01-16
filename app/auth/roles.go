package auth

// ---------------------------------------------------------------------
// Imports
// ---------------------------------------------------------------------

import (
	"scifi-search/app/infra/auth/supertokens"
)

// ---------------------------------------------------------------------
// Structures
// ---------------------------------------------------------------------

// Role in the system, with its name and auth level.
type Role struct {
	Name  string
	Level int
}

// ---------------------------------------------------------------------
//  Variables
// ---------------------------------------------------------------------

// Predefined roles.
var (
	NoRole     = Role{Name: "no-role", Level: -1}
	AdminRole  = Role{Name: "admin", Level: 2}
	LoaderRole = Role{Name: "loader", Level: 1}
	UserRole   = Role{Name: "user", Level: 0}
)

// ---------------------------------------------------------------------
//  Functions
// ---------------------------------------------------------------------

func AssignRoleToUser(role Role, userID string) {

	supertokens.AssignRoleToUser(userID, role.Name)
}

// ---------------------------------------------------------------------
