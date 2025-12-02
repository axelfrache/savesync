# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install swag CLI and generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    /go/bin/swag init -g internal/infra/http/router.go -o docs

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o savesyncd ./cmd/savesyncd

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/savesyncd .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Set environment variables
ENV SERVER_PORT=8080
ENV DATABASE_PATH=/app/data/savesync.db
ENV LOG_LEVEL=info

# Run the application
CMD ["./savesyncd"]
