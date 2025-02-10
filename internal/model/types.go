package model

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MessageTask struct {
	Message *tgbotapi.Message
	Bot     *tgbotapi.BotAPI
}

type Handler interface {
	Handler(task MessageTask) error
	CanHandle(message *tgbotapi.Message) bool
}
