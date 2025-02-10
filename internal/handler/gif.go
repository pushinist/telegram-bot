package handler

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pushinist/telegram-bot/internal/model"
)

type GifHandler struct {
	triggers map[string]string
}

func NewGifHandler() *GifHandler {
	return &GifHandler{
		triggers: map[string]string{
			"AgADzQIAAmRjXFM": "./assets/gifs/devs.gif",
		},
	}
}

func (h *GifHandler) CanHandle(message *tgbotapi.Message) bool {
	return message.Animation != nil
}

func (h *GifHandler) Handle(task model.MessageTask) error {
	if responsePath, exists := h.triggers[task.Message.Animation.FileUniqueID]; exists {
		sendGifResponse(task.Bot, task.Message, responsePath)
		return nil
	}
	return errors.New("trigger not found")
}
