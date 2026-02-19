#!/bin/bash
set -e

echo "Running DB init loader..."

run_sql_dir() {
    DIR=$1
    if [ -d "$DIR" ]; then
        echo "Running SQL files from $DIR"
        for f in "$DIR"/*.sql; do
            [ -e "$f" ] || continue
            echo "Executing $f"
            psql -v ON_ERROR_STOP=1 \
                -U "$POSTGRES_USER" \
                -d "$POSTGRES_DB" \
                -f "$f"
        done
    fi
}

if [ "$APP_ENV" = "production" ]; then
    echo "Production mode: running production SQLs"
    run_sql_dir /docker-entrypoint-initdb.d/schema
elif [ "$APP_ENV" = "development" ]; then
    echo "Development mode: running development SQLs"
    run_sql_dir /docker-entrypoint-initdb.d/init/development
else
    echo "Unknown APP_ENV='$APP_ENV', skipping environment-specific SQLs"
fi

run_sql_dir /docker-entrypoint-initdb.d/init/production

echo "DB init loader finished."
