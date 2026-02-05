package supertokens

// ------------------------------------------------------------------------------------------------
// Imports
// ------------------------------------------------------------------------------------------------

import "github.com/supertokens/supertokens-golang/recipe/userroles"

// ------------------------------------------------------------------------------------------------
// Services
// ------------------------------------------------------------------------------------------------

func CreateNewRoleOrAddPermissions(roleName string, permissions []string) {
	userroles.CreateNewRoleOrAddPermissions(roleName, permissions, nil)
}

// ------------------------------------------------------------------------------------------------

func GetRolesForUser(userID string) []string {

	roles, err := userroles.GetRolesForUser("public", userID, nil)
	if err != nil {
		return []string{}
	} else {
		return roles.OK.Roles
	}
}

// ------------------------------------------------------------------------------------------------

func AssignRoleToUser(authID, role string) {
	userroles.AddRoleToUser("public", authID, role, nil)
}

// ------------------------------------------------------------------------------------------------
