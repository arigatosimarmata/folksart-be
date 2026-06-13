# ==========================================
# Phase 1: High Performance Build Stage
# ==========================================
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Enable CGO_ENABLED=0 for absolute standalone binary compilation
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy dependency markers
COPY go.mod ./
COPY go.sum ./

# Fetch dependencies
RUN go mod download

# Copy source trees
COPY . .

# Compile and optimize binary stripping debug vectors
RUN go build -ldflags="-s -w" -o iam-governance-service main.go

# ==========================================
# Phase 2: Ultra Minimal Runtime Image
# ==========================================
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy output artifact from compiler stage
COPY --from=builder /app/iam-governance-service .

# Configure default container environment variables
ENV PORT=8080 \
    DB_USER=root \
    DB_PASSWORD=secret \
    DB_HOST=mysql-db \
    DB_PORT=3306 \
    DB_NAME=folksart_iam

EXPOSE 8080

CMD ["./iam-governance-service"]
