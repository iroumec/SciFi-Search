#!/bin/bash

# This script is used to generate type-safe code for database queries using sqlc.
# It copies the sqlc configuration file from the resources directory, runs the sqlc generate command,
# and then removes the copied configuration file.

cp database/sqlc.yaml .
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
sqlc generate
rm sqlc.yaml