# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies and HTTPS)
RUN apk add --no-cache git ca-certificates tzdata

# Create appuser for security
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Ensure static directory exists
RUN mkdir -p static

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o whatsmeow-mcp .

# Final stage
FROM scratch

# Import ca-certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Import user and group files from builder
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /build/whatsmeow-mcp /whatsmeow-mcp

# Copy migrations
COPY --from=builder /build/migrations /migrations

# Copy static files
COPY --from=builder /build/static /static

# Use non-root user
USER appuser

# Expose port (adjust if your app uses a different port)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/whatsmeow-mcp", "--health-check"] || exit 1

# Run the binary
ENTRYPOINT ["/whatsmeow-mcp"]
