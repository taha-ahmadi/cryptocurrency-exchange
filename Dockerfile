FROM golang:1.20.1-alpine3.16 AS build

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app

# Final stage
FROM alpine:latest

RUN apk add --no-cache bash

# Copy the built binary from the build stage
COPY --from=build . .

EXPOSE 3000

# Set the working directory
WORKDIR /app

# Run the application
CMD ["/bin/app"]


