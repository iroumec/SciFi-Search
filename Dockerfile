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
RUN go build -buildvcs=false -o /app/main ./app

# ============================================================
# Stage 2: Development
# ============================================================

FROM golang:1.25-alpine AS dev

WORKDIR /app

# Se instalan "git" y "bash".
RUN apk add --no-cache git bash

# Se instalan dependencias necesarias para la validación de la constancia de alumno regular.
RUN apk add --no-cache poppler-utils chromium nss freetype harfbuzz ttf-freefont

# Se instala Air.
RUN go install github.com/air-verse/air@latest

# Se expone el puerto 8080.
EXPOSE 8080

# Ejecución de Air.
CMD air -c .air.toml

# ============================================================
# Stage 3: Production
# ============================================================

FROM alpine:latest AS production

WORKDIR /app

# Se instalan dependencias necesarias para la validación de la constancia de alumno regular.
RUN apk add --no-cache poppler-utils chromium nss freetype harfbuzz ttf-freefont

# Se crea un usuario y un grupo.
RUN addgroup -S olimpiadas-unicen && adduser -S olimpiadas-unicen -G olimpiadas-unicen

# Se copia en la imagen el binario y los assets desde el builder.
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/template ./template


RUN chown -R olimpiadas-unicen:olimpiadas-unicen /app

EXPOSE 8080
USER olimpiadas-unicen
CMD ["./main"]