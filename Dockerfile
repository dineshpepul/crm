Use an official Golang runtime as a parent image
FROM golang:alpine
LABEL authors="dinesh"

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application
RUN go build -o main .

# Expose a port (if your Go application listens on a specific port)
EXPOSE 8005

# Run the Go application when the container starts
CMD ["./main"]