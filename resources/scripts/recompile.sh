#!/bin/bash

# Detiene el script si algo falla
set -e

echo "Recompilando todo..."

# Generación de sqlc.
./resources/scripts/sqlc/runSQLC.sh

# Generación de templ.
./resources/scripts/templ/runTEMPL.sh

# Compilación de Go
go build -buildvcs=false -o ./tmp/main ./app