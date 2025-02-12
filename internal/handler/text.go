package handler

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pushinist/telegram-bot/internal/model"
)

type TextHandler struct {
	triggers map[string]string
}

func NewTextHandler() *TextHandler {
	return &TextHandler{
		triggers: map[string]string{
			"разрабы дауны": "./assets/gifs/devs.gif",
			"разраб даун":   "./assets/gifs/devs.gif",
			"даун":          "./assets/gifs/devs.gif",
			"йоу":           "./assets/gifs/cat-wif.gif",
		},
	}
}

func (h *TextHandler) CanHandle(message *tgbotapi.Message) bool {
	return message.Text != ""
}

func (h *TextHandler) Handle(task *model.MessageTask) error {
	for trigger, gifPath := range h.triggers {
		if strings.Contains(trigger, strings.ToLower(task.Message.Text)) {
			return sendGifResponse(task.Bot, task.Message, gifPath)
		}
	}
	return nil
}
