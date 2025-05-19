FROM golang:1.24.3-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Copy the vendor directory to avoid downloading dependencies again
COPY vendor/ ./vendor/

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/server ./src/cmd/server

# Use a minimal alpine image for the final image
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bin/server /app/server

# Copy docs directory (don't regenerate)
COPY --from=builder /app/docs /app/docs

# Set environment variables
ENV GIN_MODE=release

# Expose the application port
EXPOSE 8585

# Run the application
CMD ["/app/server"]
