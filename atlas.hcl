# atlas.hcl

# Variables de entorno.
variable "db_url" {
  type    = string
  default = getenv("DB_URL")
}

variable "atlas_db_url" {
  type    = string
  default = getenv("ATLAS_DB_URL")
}

# Docker enviroment (inside container).
env "docker" {

  # Main database URL.
  url = "${var.db_url}"

  # Atlas database.
  # It is required so Atlas can calculate temporal differences in the scheme.
  # An extension used for trigrams is specified.
  dev = "${var.atlas_db_url}"
  
  # Migrations directory.
  migration {
    dir = "file://database/migrations"
  }
  
  # Source scheme (what it is edited).
  # Queries don't require migrations.
  src = "file://database/schema/schema.sql"
}

# Local environment (outside contenedor).
env "local" {

  url = "${var.db_url}"

  dev = "${var.atlas_db_url}"
  
  migration {
    dir = "file://database/migrations"
  }
  
  src = "file://database/schema/schema.sql"
}