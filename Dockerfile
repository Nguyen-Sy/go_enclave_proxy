# Build stage
FROM golang:1.23-alpine AS builder

# Install required build tools
RUN apk add --no-cache gcc musl-dev linux-headers

# Set working directory
WORKDIR /app

# Copy Go module files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN go build -o ec2-client

# Final stage
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates

# Copy the built executable
COPY --from=builder /app/ec2-client /usr/local/bin/ec2-client

# Create a non-root user
RUN adduser -D appuser
USER appuser

# Set the entry point
ENTRYPOINT ["ec2-client"]