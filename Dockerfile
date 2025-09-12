# Etapa de build
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Copiar archivos de módulos
COPY go.mod go.sum ./
RUN go mod tidy

# Copiar todo el código Go
COPY app ./app

# Compilar la app
WORKDIR /app/app
RUN go build -o /app/main .

# Imagen final ligera
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
