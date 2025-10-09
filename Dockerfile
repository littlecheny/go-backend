# syntax=docker/dockerfile:1

# ---- Builder stage ----
FROM golang:1.24 AS builder
WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy source
COPY . .

# Build (static)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/main.go

# ---- Runtime stage ----
FROM alpine:3.20
# Add non-root user
RUN adduser -D -h /home/app appuser
WORKDIR /home/app

# Copy binary
COPY --from=builder /app/server /usr/local/bin/server

# Runtime envs
ENV GIN_MODE=release

# Expose default port (match your .env, e.g. :8080)
EXPOSE 8080

# Ensure log file path is writable
RUN mkdir -p /home/app && chown -R appuser:appuser /home/app

USER appuser
ENTRYPOINT ["/usr/local/bin/server"]