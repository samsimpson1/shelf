# Build stage
FROM golang:alpine AS builder

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies (if go.sum exists)
RUN go mod download

# Copy source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shelf .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests to TMDB API
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/shelf .

# Copy templates and static files
COPY templates/ ./templates/
COPY static/ ./static/

# Expose port 8080
EXPOSE 8080

# Environment variables with defaults
ENV PORT=8080
ENV MEDIA_DIR=/media
ENV IMPORT_DIR=""
ENV TMDB_API_KEY=""

# Run the application
CMD ["./shelf"]
