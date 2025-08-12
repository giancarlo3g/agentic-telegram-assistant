# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with debug information
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -gcflags="all=-N -l" -o calendar-bot ./cmd/bot

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates and debugging tools
RUN apk --no-cache add ca-certificates curl busybox-extras

# Copy the binary from builder stage
COPY --from=builder /app/calendar-bot .

# Copy environment example
COPY --from=builder /app/env.example .

# Create directories for logs and credentials
RUN mkdir -p /app/logs /app/credentials

# Expose port
EXPOSE 8080

# Run the application
CMD ["./calendar-bot"]
