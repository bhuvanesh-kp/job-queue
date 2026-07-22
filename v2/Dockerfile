# Use the official Golang image for building
FROM golang:1.25 AS builder

# Set working directory
WORKDIR /app

# Copy Go modules and dependencies
COPY go.mod ./
RUN go mod download

# Copy source code
COPY main.go .

# Build the application
RUN go build -o /main 

##################### multistage docker build ##########################

# Use a minimal base image for final deployment
# FROM alpine:latest

# # Set working directory in the container
# WORKDIR /root/

# # Copy the built binary from the builder stage
# COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]