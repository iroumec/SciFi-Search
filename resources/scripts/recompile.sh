#!/bin/bash

# Detiene el script si algo falla
set -e

echo
echo "=================================================================================="
echo "Compilando todo... ¡Espere, por favor! Esto puede tomar un tiempo la primera vez."
echo "=================================================================================="
echo

# Generación de sqlc.
./resources/scripts/sqlc/runSQLC.sh

# Generación de templ.
./resources/scripts/templ/runTEMPL.sh

# Compilación de Go
go build -buildvcs=false -o ./tmp/main ./app

echo
echo "=================================================================================="
echo "¡Compilación finalizada!"
echo "=================================================================================="
echo