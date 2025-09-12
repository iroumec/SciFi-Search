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
COPY --from=builder /app/main .
COPY static ./static
EXPOSE 8080
CMD ["./main"]
