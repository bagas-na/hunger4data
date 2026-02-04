### Build stage 
FROM golang:1.25.6-alpine AS builder

WORKDIR /workspace

# Dibutuhkan untuk TLS dan HTTPS
RUN apk add --no-cache git ca-certificates

# Go workspace
COPY go.work go.work.sum ./

# Copy proto file
COPY proto/ proto/

# Copy semua microservice (dibutuhkan karena perlu reference ke semua jika menggunakan go.work)
COPY apps ./apps

WORKDIR /workspace/apps/payment

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/payment-service ./cmd

### Runtime stage
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/payment-service /app/payment-service

# gRPC + HTTP webhook
EXPOSE 50054
EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/app/payment-service"]
