#!/bin/bash
set -euo pipefail

for file in database/post-migrate/*.sql; do
    psql "$DB_URL" -v ON_ERROR_STOP=1 -f "$file"
done
