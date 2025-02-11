FROM golang:1.23.6-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o telegram-bot ./cmd/bot/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/telegram-bot .
COPY asstes/gifs ./assets/gifs

CMD ["./telegram-bot"]
