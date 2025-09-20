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

# Compilamos el paquete que está dentro de la carpeta ./app
RUN go build -buildvcs=false -o /app/main ./app

# ============================================================
# Stage 2: Development
# ============================================================
FROM golang:1.25-alpine AS dev
WORKDIR /app
RUN apk add --no-cache git bash
RUN go install github.com/air-verse/air@latest
ENV PATH=$PATH:/root/go/bin
EXPOSE 8080
# Usamos el formato shell que ya nos funcionó
CMD air -c .air.toml

# ============================================================
# Stage 3: Production
# ============================================================
FROM alpine:latest AS production
WORKDIR /app
RUN apk add --no-cache poppler-utils chromium nss freetype harfbuzz ttf-freefont
RUN addgroup -S uki && adduser -S uki -G uki

# Copiamos el binario y los assets desde el builder
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/template ./template

RUN chown -R uki:uki /app
EXPOSE 8080
USER uki
CMD ["./main"]