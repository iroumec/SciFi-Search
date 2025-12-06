# ============================================================
# Stage 1: Builder
# ============================================================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copiamos el ÚNICO go.mod/go.sum desde la raíz
COPY go.mod go.sum ./
RUN go mod download

# Copiamos TODO el código fuente
COPY . .

# Se compila el paquete dentro de ./app.
# La opción "-buildvcs=false" la pedía una dependencia.
# Se usa CGO_ENABLED=0 para binario estático y output correcto.
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o main ./app

# ============================================================
# Stage 2: Assets (Tailwind)
# ============================================================

# Instalación de Node.
FROM node:20-alpine AS assets

WORKDIR /assets

# Se copian los archivos de configuración de npm.
COPY package.json ./

# Se copia el archivo de configuración de Tailwind.
COPY tailwind.config.js ./

# Copiado de templates.
COPY static ./static
COPY app/views ./app/views

# Instalación de dependencias.
# Si se instalan desde el lock, usar "ci".
# Se genera el package-lock.json la primera vez que se ejecute npm install.
# En builds posteriores, si el archivo existe, lo usará.
RUN npm install && npm install tailwindcss postcss autoprefixer

# Compilación de Tailwind.
RUN npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify

# ============================================================
# Stage 3: Development
# ============================================================

FROM golang:1.25-alpine AS dev

WORKDIR /app

# Se instalan las herramientas base: "git", "bash" y "curl".
RUN apk add --no-cache git bash curl

# Instalación de Node (para Tailwind).
RUN apk add --no-cache nodejs npm

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

# Copiado de las dependencias de Tailwind.
COPY package.json tailwind.config.js ./

# Instalación de dependencias.
RUN npm install

# Se expone el puerto 8080.
EXPOSE 8080

# Ejecución de Air.
CMD air -c .air.toml

# ============================================================
# Stage 4: Production
# ============================================================

FROM alpine:latest AS production

WORKDIR /app

# Se crea un usuario y un grupo.
RUN addgroup -S sciFi-search && adduser -S sciFi-search -G sciFi-search

# Se copia en la imagen el binario y los assets desde el builder.
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/resources/planillas ./resources/planillas

# Se copia el CSS generado por Tailwind.
COPY --from=assets /assets/static/css/output.css ./static/css/

RUN chown -R sciFi-search:sciFi-search /app && \
    chmod +x /app/main

EXPOSE 8080
USER sciFi-search
CMD ["./main"]