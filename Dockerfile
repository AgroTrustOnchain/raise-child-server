# ─────────────────────────────────────────────────────────────
# Stage 1: Builder — full Go toolchain, only used at build time
# ─────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Cache go modules layer separately — only re-downloaded when go.sum changes
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a fully static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/server .

# ─────────────────────────────────────────────────────────────
# Stage 2: Runtime — scratch (empty) base for minimal image
# ─────────────────────────────────────────────────────────────
FROM scratch

# Copy TLS root certificates so HTTPS calls work inside scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary only
COPY --from=builder /app/server /server

EXPOSE 8080
ENTRYPOINT ["/server"]