# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /booking-api ./cmd/api

# Final runner stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /booking-api /app/booking-api

EXPOSE 8080

CMD ["/app/booking-api"]