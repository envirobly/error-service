FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

# Build the Go app for the Linux arm64 architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o error-service

# Stage 2: Create the minimal image
FROM scratch

# Copy the built Go binary from the builder stage
COPY --from=builder /app/error-service /error-service

# Set the entrypoint to the built Go binary
ENTRYPOINT ["/error-service"]
