package handler

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pushinist/telegram-bot/internal/model"
	"log/slog"
	"strings"
)

type TextHandler struct {
	triggers map[string]string
}

func NewTextHandler() *TextHandler {
	return &TextHandler{
		triggers: map[string]string{
			"разрабы дауны": "./gifs/devs.gif",
			"разраб даун":   "./gifs/devs.gif",
			"даун":          "./gifs/devs.gif",
		},
	}
}

func (h *TextHandler) CanHandle(message *tgbotapi.Message) bool {
	return message.Text != ""
}

func (h *TextHandler) Handle(task model.MessageTask) error {
	for trigger, gifPath := range h.triggers {
		if strings.Contains(trigger, strings.ToLower(task.Message.Text)) {
			slog.Info("Trying to send")
			sendGifResponse(task.Bot, task.Message, gifPath)
			return nil
		}
	}
	return errors.New("trigger not found")
}
