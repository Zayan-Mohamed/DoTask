# Multi-stage build for Go backend
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install necessary packages
RUN apk add --no-cache git wget

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Create .env file with Docker-specific configuration
RUN echo 'DATABASE_URL=postgres://postgres:admin@postgres:5432/dotask?sslmode=disable' > .env && \
    echo 'PORT=8080' >> .env && \
    echo 'JWT_SECRET=81b705c0597375a27d6d08b8b09b9dd073acf15218128970973477e429c4870b' >> .env && \
    echo 'CORS_ORIGINS=http://localhost,http://localhost:80' >> .env && \
    echo 'GIN_MODE=release' >> .env

EXPOSE 8080

CMD ["./main"]
