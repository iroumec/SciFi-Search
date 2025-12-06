#!/bin/bash

# Detiene el script si algo falla
set -e

echo
echo "=================================================================================="
echo "Compilando todo... ¡Espere, por favor! Esto puede tomar un tiempo la primera vez."
echo "=================================================================================="
echo

# 1. Generaración y aplicación migraciones con Atlas.
echo "Verificando migraciones de base de datos..."

# Generación de migraciones.
atlas migrate diff --env docker

# Aplicación de migraciones.
atlas migrate apply --env docker

# 2. Generación de sqlc.
echo "Generando sqlc..."

# Copia del archivo .yaml.
cp database/sqlc.yaml .

# Ejecución de SQLC.
# Se puede ejecutar usando solo "sqlc"
# ya que se exportó al PATH en el DockerFile.
sqlc generate

# Borrado del yaml.
rm sqlc.yaml

# 3. Generación de templ
echo "Generando templ..."
go run github.com/a-h/templ/cmd/templ@latest generate ./app/views

# 4. Compilación de Tailwind...
echo "Ejecutando TailwindCSS..."
npm run build:css

# 5. Compilación de Go
echo "Compilando aplicación Go..."
go build -buildvcs=false -o ./tmp/main ./app

echo
echo "=================================================================================="
echo "¡Compilación finalizada!"
echo "=================================================================================="
echo