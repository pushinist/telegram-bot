FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o telegram-bot ./cmd/bot/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/telegram-bot .
COPY --from=builder /app/.env .
COPY assets/gifs ./assets/gifs

CMD ["./telegram-bot"]