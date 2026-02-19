# =================================================================================================
# Stage 1: Builder
# =================================================================================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Se copia go.mod/go.sum.
COPY go.mod go.sum ./
RUN go mod download

# Se copia todo el código fuente
COPY . .

# Se compila el paquete dentro de ./app.
RUN CGO_ENABLED=0 GOOS=linux \
    # Se usa CGO_ENABLED=0 para binario estático y output correcto.
    go build \
    -trimpath \
    -ldflags="-s -w" \
    # El anterior comando reduce los símbolos de debug (achicando así la imagen).
    -buildvcs=false \
    # Al anterior comando lo pedía una dependencia.
    -o main ./app

# =================================================================================================
# Stage 2: Development
# =================================================================================================

FROM golang:1.25-alpine AS dev

WORKDIR /app

# Se instalan las herramientas base: "git", "bash" y "curl".
RUN apk add --no-cache git bash curl

ENV PATH="/go/bin:${PATH}"

# Se instalan Air, SQLC y Templ.
RUN go install github.com/air-verse/air@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    go install github.com/a-h/templ/cmd/templ@latest

# Se instala Atlas CLI.
RUN curl -sSfL https://release.ariga.io/atlas/atlas-linux-amd64-latest \
    -o /usr/local/bin/atlas && chmod +x /usr/local/bin/atlas

# Se verifica la instalación de Atlas.
RUN atlas version

# Se expone el puerto 8080.
EXPOSE 8080

# Ejecución de Air.
CMD air -c .air.toml

# =================================================================================================
# Stage 3: Production
# =================================================================================================

FROM alpine:3.20 AS production

WORKDIR /app

# Se crea un usuario y un grupo.
RUN addgroup -S sciFi-search && adduser -S sciFi-search -G sciFi-search

# Se copia en la imagen el binario y los recursos requeridos desde el builder.
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/resources/planillas ./resources/planillas
COPY --from=builder /app/resources/languages ./resources/languages

RUN chown -R sciFi-search:sciFi-search /app && \
    chmod +x /app/main

# Se expone el puerto 8080.
EXPOSE 8080

# Se cambia el usuario al definido.
USER sciFi-search

# Se ejecuta la aplicación.
CMD ["./main"]

# =================================================================================================