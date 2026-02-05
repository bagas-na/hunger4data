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

WORKDIR /workspace/apps/api-gateway

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/api-gateway ./cmd

### Runtime stage
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/api-gateway /app/api-gateway

EXPOSE 50051

USER nonroot:nonroot

ENTRYPOINT ["/app/api-gateway"]
