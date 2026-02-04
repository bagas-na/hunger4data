### Build Stage ###
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

WORKDIR /workspace/apps/notification

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/notification ./cmd

### Runtime stage ###
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/notification /app/notification

EXPOSE 50053

USER nonroot:nonroot

ENTRYPOINT ["/app/notification"]
