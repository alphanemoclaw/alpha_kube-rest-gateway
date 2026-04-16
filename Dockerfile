# --- Build Stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for a static binary that works in a minimal image
RUN CGO_ENABLED=0 GOOS=linux go build -o kube-rest-gateway main.go

# --- Final Stage ---
FROM alpine:3.19

WORKDIR /app

# Add a non-root user for security
RUN adduser -D gatewayuser
USER gatewayuser

# Copy the binary from the builder stage
COPY --from=builder /app/kube-rest-gateway .

# Expose the default port
EXPOSE 8081

# Set default environment variables
ENV GATEWAY_PORT=8081
ENV GATEWAY_HOST=0.0.0.0
ENV LOG_LEVEL=INFO

# Command to run the gateway
ENTRYPOINT ["./kube-rest-gateway"]
