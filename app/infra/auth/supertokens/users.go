package supertokens

import "github.com/supertokens/supertokens-golang/supertokens"

func DeleteUser(authID string) error {
	err := supertokens.DeleteUser(authID)

	if err != nil {
		return err
	}

	// Usuario eliminado exitosamente.
	return nil
}
