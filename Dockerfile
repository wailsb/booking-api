# Change from golang:1.23-alpine to golang:1.27-alpine (or latest)
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /booking-api ./cmd/server
# Final lightweight image stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /booking-api .
EXPOSE 8080
CMD ["./booking-api"]