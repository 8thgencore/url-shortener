# Use the official Golang image as the base image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module and Go sum files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the Go application
RUN go build -o url-shortener ./cmd/url-shortener

# Expose the port that the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./url-shortener"]
