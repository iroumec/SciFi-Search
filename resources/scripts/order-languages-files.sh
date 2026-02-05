#!/bin/bash

find "resources/languages" -type f -name "*.json" -print0 | while IFS= read -r -d '' file; do
    jq -S --indent 4 . "$file" > "${file}.tmp" && mv "${file}.tmp" "$file"
done
