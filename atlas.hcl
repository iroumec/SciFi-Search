# atlas.hcl

# Variables de entorno.
variable "db_user" {
  type    = string
  default = getenv("DB_USER")
}

variable "db_password" {
  type    = string
  default = getenv("DB_PASSWORD")
}

variable "db_name" {
  type    = string
  default = getenv("DB_NAME")
}

variable "db_port" {
  type    = string
  default = getenv("DB_PORT")
}

variable "db_host" {
  type    = string
  default = getenv("DB_HOST")
}

# Entorno Docker (dentro del contenedor).
env "docker" {

  # URL de la base de datos principal.
  url = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=disable"

  # Atlas requiere que se le especifique una base de datos
  # temporal para calcular diferencias en el schema.
  dev = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/atlas_temp?sslmode=disable"
  
  # Directorio donde se guardan las migraciones.
  migration {
    dir = "file://database/migrations"
  }
  
  # Schema fuente (lo que se edita).
  # Las queries no requieren de migración.
  src = "file://database/schema/schema.sql"
}

# Entorno local (fuera del contenedor).
env "local" {

  url = "postgres://${var.db_user}:${var.db_password}@localhost:${var.db_port}/${var.db_name}?sslmode=disable"

  # Atlas requiere que se le especifique una base de datos
  # temporal para calcular diferencias en el schema.
  dev = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/atlas_temp?sslmode=disable"
  
  migration {
    dir = "file://database/migrations"
  }
  
  src = "file://database/schema/schema.sql"
}