#!/bin/bash
set -euo pipefail

BASE_DIR="resources/languages"

declare -A COUNTS

# --------------------------------------------------
# Contar traducciones por idioma
# --------------------------------------------------
for lang_dir in "$BASE_DIR"/*; do
    [ -d "$lang_dir" ] || continue
    lang="$(basename "$lang_dir")"

    count=0
    while IFS= read -r -d '' file; do
        file_count=$(jq '[paths(scalars)] | length' "$file")
        count=$((count + file_count))
    done < <(find "$lang_dir" -type f -name "*.json" -print0)

    COUNTS["$lang"]=$count
done

# --------------------------------------------------
# Comparación
# --------------------------------------------------
expected=""
for lang in "${!COUNTS[@]}"; do
    if [[ -z "$expected" ]]; then
        expected="${COUNTS[$lang]}"
    elif [[ "${COUNTS[$lang]}" -ne "$expected" ]]; then
        echo "ERROR: Cantidad de traducciones inconsistente."
        for l in "${!COUNTS[@]}"; do
            echo "  $l: ${COUNTS[$l]}"
        done
        exit 1
    fi
done

echo "OK: Todos los idiomas tienen $expected traducciones."
