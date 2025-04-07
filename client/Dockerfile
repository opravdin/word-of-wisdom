FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o word-of-wisdom-client ./cmd/client

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/word-of-wisdom-client .
COPY --from=builder /app/internal/http/static ./internal/http/static

EXPOSE 3000
CMD ["./word-of-wisdom-client", "-server", "server:8080"]