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
# La opción "-buildvcs=false" la pedía una dependencia.
# Se usa CGO_ENABLED=0 para binario estático y output correcto.
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o main ./app

# =================================================================================================
# Stage 2: Development
# =================================================================================================

FROM golang:1.25-alpine AS dev

WORKDIR /app

# Se instalan las herramientas base: "git", "bash" y "curl".
RUN apk add --no-cache git bash curl

# Se instala Air.
RUN go install github.com/air-verse/air@latest

# Se instala SQLC y se exporta al PATH.
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN export PATH="$PATH:$(go env GOPATH)/bin"

# Se instala Templ.
RUN go install github.com/a-h/templ/cmd/templ@latest

# Se instala Atlas CLI.
RUN curl -sSfL https://release.ariga.io/atlas/atlas-linux-amd64-latest -o /usr/local/bin/atlas && \
    chmod +x /usr/local/bin/atlas

# Se verifica la instalación de Atlas.
RUN atlas version

# Se expone el puerto 8080.
EXPOSE 8080

# Ejecución de Air.
CMD air -c .air.toml

# =================================================================================================
# Stage 3: Production
# =================================================================================================

FROM alpine:latest AS production

WORKDIR /app

# Se crea un usuario y un grupo.
RUN addgroup -S sciFi-search && adduser -S sciFi-search -G sciFi-search

# Se copia en la imagen el binario y los assets desde el builder.
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/resources/planillas ./resources/planillas
COPY --from=builder /app/resources/languages ./resources/languages

RUN chown -R sciFi-search:sciFi-search /app && \
    chmod +x /app/main

EXPOSE 8080
USER sciFi-search
CMD ["./main"]

# =================================================================================================