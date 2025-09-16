# Builder
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go mod tidy
RUN go build -o /app/main ./app

# Imagen final
FROM alpine:latest

WORKDIR /app

# Creación de un usuario y grupo no root.
RUN addgroup -S uki && adduser -S uki -G uki

COPY --from=builder /app/main .

COPY static ./static

# Se brinda al usuario no root permisos para ejecutar la app.
RUN chown -R uki:uki /app

EXPOSE 8080

# Se cambia al usuario no root.
USER uki

CMD ["./main"]