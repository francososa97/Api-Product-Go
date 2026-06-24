# ---- Etapa de build ----
FROM golang:1.22-alpine AS builder

WORKDIR /src

# Se copian primero los manifiestos para aprovechar la caché de capas de Docker.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Binario estático (sin CGO) para poder correr en una imagen mínima.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/product-api ./cmd/api

# ---- Etapa final ----
FROM alpine:3.20

# Certificados para conexiones TLS salientes (p. ej. MongoDB Atlas).
RUN apk add --no-cache ca-certificates

# Usuario sin privilegios.
RUN adduser -D -u 10001 appuser
USER appuser

COPY --from=builder /bin/product-api /usr/local/bin/product-api

EXPOSE 8080
ENTRYPOINT ["product-api"]
