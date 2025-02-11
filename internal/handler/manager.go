package handler

import (
	"fmt"
	"github.com/pushinist/telegram-bot/internal/model"
	"log/slog"
)

type Manager struct {
	handlers []model.Handler
}

func NewManager() *Manager {
	return &Manager{
		handlers: []model.Handler{
			NewTextHandler(),
			NewGifHandler(),
		},
	}
}

func (m *Manager) Handle(task model.MessageTask) error {
	for _, handler := range m.handlers {
		if handler.CanHandle(task.Message) {
			slog.Info(fmt.Sprintf("Handler found for message %d", task.Message.MessageID))
			return handler.Handle(&task)
		}
	}
	return nil

}
