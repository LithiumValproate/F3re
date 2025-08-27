# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set the working directory for the build
WORKDIR /build

# Copy all services source code needed for the build.
# This is the simplest way to handle local `replace` directives in go.mod.
COPY ./services/ /build/services/

# Change to the specific service directory we want to build
WORKDIR /build/services/go-chat

# Build the binary. CGO_ENABLED=0 is used for creating a static binary.
# The output is placed at /app/server inside the builder stage.
RUN CGO_ENABLED=0 go build -o /app/server .

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/server .

# Expose port 8080 from the container
EXPOSE 8080

# The command to run when the container starts
CMD ["./server"]
