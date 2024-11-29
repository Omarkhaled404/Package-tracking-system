# Use the official Golang image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /build

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go build-o./userapi./cad/userapi/main.go
# Copy the entire project into the container
COPY . .
# Build the Go application
RUN go build -o bin main.go


FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY -from-builder/build/userapi./userapi
CMD ["/app/userapi"]
