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

atlas migrate hash --dir file://database/migrations

# Generación de migraciones.
atlas migrate diff --env docker

# Aplicación de migraciones.
atlas migrate apply --env docker

# 2. Generación de sqlc.
echo "Generando sqlc..."
sqlc generate -f database/sqlc.yaml

# 3. Generación de templ
echo "Generando templ..."
go run github.com/a-h/templ/cmd/templ@latest generate ./app/views

# 4. Compilación de Go
echo "Compilando aplicación Go..."
go build -buildvcs=false -o ./tmp/main ./app

echo
echo "=================================================================================="
echo "¡Compilación finalizada!"
echo "=================================================================================="
echo