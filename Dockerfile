# Use a lightweight Go image
FROM golang:1.22 AS builder

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies and build the app
RUN go mod tidy && go build -o url-shortener ./cmd/server/main.go

# Use a minimal image for final container
FROM alpine:latest

# Set working directory in the final container
WORKDIR /app

# Copy built binary from builder stage
COPY --from=builder /app/url-shortener .

# Expose port 8080 for the app
EXPOSE 8080

# Run the application
CMD ["./url-shortener"]
