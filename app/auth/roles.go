package auth

// Role representa un rol del sistema con su nombre y nivel de autorización.
type Role struct {
	Name  string
	Level int
}

// Roles predefinidos del sistema.
var (
	NoRole     = Role{Name: "no-role", Level: -1}
	AdminRole  = Role{Name: "admin", Level: 2}
	LoaderRole = Role{Name: "loader", Level: 1}
	UserRole   = Role{Name: "user", Level: 0}
)
