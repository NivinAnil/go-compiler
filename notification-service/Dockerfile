FROM golang:1.23-alpine3.19
# Add Maintainer Info
VOLUME /var/log/
ENV LOG_PATH=/var/log/
ARG environment
ENV BOOT_CUR_ENV=${environment}

# Set the Current Working Directory inside the container
WORKDIR /app

RUN apk update
RUN apk add --no-cache git

ENV GOPRIVATE="github.com/kodnest"
RUN git config --global url."https://VikramRay96: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@github.com".insteadOf "https://github.com"
# Copy go mod and sum files
COPY go.mod .
COPY go.sum .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY ./common ./common
COPY ./notification-service ./notification-service

# Build the Go app
RUN go build -o main ./notification-service

# Expose port 8080 to the outside world
EXPOSE 8080

CMD  ["sh", "-c", "./main"]