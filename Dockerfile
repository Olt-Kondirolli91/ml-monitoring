# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the files
COPY . .

# Build the application
RUN go build -o /ml-monitoring ./cmd/main.go

# Final stage
FROM alpine:3.21
WORKDIR /root/
COPY --from=builder /ml-monitoring ./
COPY migrations ./migrations
CMD ["./ml-monitoring"]
