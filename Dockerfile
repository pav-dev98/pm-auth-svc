# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Descarga de dependencias (capa cacheable)
COPY go.mod go.sum ./
RUN go mod download

# Código fuente
COPY . .

# Binario estático para correr en una imagen mínima
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /auth-svc ./cmd

# ---- Runtime stage ----
FROM alpine:3.20

# Certificados para conexiones TLS salientes
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Usuario sin privilegios
RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=builder /auth-svc /app/auth-svc

EXPOSE 50051

ENTRYPOINT ["/app/auth-svc"]
