FROM golang:1.23.6-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o telegram-bot ./cmd/bot/main.go

FROM alpine:latest

WORKDIR /app
RUN mkdir -p ./assets/gifs
COPY --from=builder /app/telegram-bot .
COPY assets/gifs ./assets/gifs

RUN chmod +x telegram-bot
CMD ["./telegram-bot"]
