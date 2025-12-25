## DEVELOPMENT STAGE
FROM golang:1.25-trixie AS development

# Install git for module downloads
RUN apt-get update && \
    apt-get install -y git

# Set direktori kerja
WORKDIR /app

# Copy source code
COPY . .

# Copy go mod files
# COPY go.work go.work.sum ./

# Download dependencies
RUN go work sync

# Install Watch tool untuk live reload saat development
RUN go install github.com/air-verse/air@latest

# Instal Delve untuk debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Build Producer
RUN go build -o /app/main /app/webcore/main.go

EXPOSE 7272

# CMD ["air", "--build.cmd", "go build -o /app/consumer /app/cmd/consumer/main.go", "--build.bin", "/app/consumer", "--debug.host", "0.0.0.0", "--debug.port", "2345"]
CMD ["air"]

## BUILD STAGE
FROM golang:1.25-trixie AS builder

# Install git for module downloads
RUN apt-get update && \
    apt-get install -y git

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Copy go mod files
# COPY go.work go.work.sum ./

# Download dependencies
RUN go work sync

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main /app/webcore/main.go

## PRODUCTION STAGE
FROM debian:trixie-slim AS production

# Set working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy configuration files
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/access.yaml .

# Expose port
EXPOSE 7272

# Run the application
CMD ["./main"]
