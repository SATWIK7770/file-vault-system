# Use official Go image as builder
FROM golang:1.25.6 AS builder

# Set working directory
WORKDIR /app

# Install git and other dependencies
RUN apk add --no-cache git

# Copy go modules files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o fvs-backend ./cmd/main.go

# Second stage: small runtime image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/fvs-backend .

# Copy any static assets (optional)
# COPY --from=builder /app/uploads ./uploads

# Expose port
EXPOSE 8080

# Command to run
CMD ["./fvs-backend"]
