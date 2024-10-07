# Use an official Golang image to build the worker binary
FROM golang:1.18-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go worker application
RUN go build -o worker .

# Use an official Python runtime as the final image for executing Python code
FROM python:3.9-alpine

# Set the working directory inside the container
WORKDIR /app

# Install necessary system dependencies, including Python 3 and bash
RUN apk add --no-cache bash
RUN apk add --no-cache nodejs npm

# Copy the built Go binary from the builder stage
COPY --from=builder /frontend/worker /app/worker

# Set the environment variables for RabbitMQ
ENV RABBITMQ_URL amqp://guest:guest@localhost:5672/

# Expose port 5672 for RabbitMQ
EXPOSE 5672

# Run the worker binary
CMD ["/app/worker"]
