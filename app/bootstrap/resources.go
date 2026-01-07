package bootstrap

import (
	"database/sql"

	"scifi-search/app/database"
)

type Resources struct {
	DB      *sql.DB
	Queries *database.Queries
}
