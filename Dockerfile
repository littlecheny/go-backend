# ---- Builder stage ----
FROM golang:1.24 AS builder
WORKDIR /app

# 设置Go代理为国内镜像源（阿里云）
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/,https://goproxy.cn,direct
ENV GO111MODULE=on

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Install swag for generating swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source
COPY . .

# Generate swagger docs
RUN /go/bin/swag init -g cmd/main.go

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