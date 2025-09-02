# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o wedding_app ./cmd/wedding_website/http

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/wedding_app .
COPY internal/templates ./internal/templates
COPY internal/static ./internal/static
COPY config ./config

EXPOSE 8081

CMD ["./wedding_app"]
