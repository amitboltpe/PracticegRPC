# Build stage
FROM golang:1.22 AS builder
WORKDIR /app

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go binary
RUN go build -o mytestapp .

# Runtime stage
FROM debian:bullseye-slim
WORKDIR /app

COPY --from=builder /app/mytestapp .

EXPOSE 9091

CMD ["./mytestapp"]