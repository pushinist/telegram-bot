package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pushinist/telegram-bot/internal/config"
	"github.com/pushinist/telegram-bot/internal/handler"
	"github.com/pushinist/telegram-bot/internal/model"
	"log/slog"

	"sync"
)

type Bot struct {
	api      *tgbotapi.BotAPI
	tasks    chan model.MessageTask
	workers  int
	wg       sync.WaitGroup
	handlers *handler.Manager
}

func New(cfg *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:      api,
		tasks:    make(chan model.MessageTask, 100),
		workers:  cfg.Workers,
		handlers: handler.NewManager(),
	}, nil
}

func (b *Bot) Start() {
	slog.Info(fmt.Sprintf("Starting %d workers", b.workers))
	for i := 1; i <= b.workers; i++ {
		b.wg.Add(1)
		go b.startWorker(i)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		slog.Info(fmt.Sprintf("Recieved message: %v", update.Message))
		select {
		case b.tasks <- model.MessageTask{
			Message: update.Message,
			Bot:     b.api,
		}:
			slog.Info(fmt.Sprintf("Processing message: %v", update.Message))
		default:
			slog.Info(fmt.Sprintf("Channel is full, skipping message: %v", update.Message))

		}

	}
}

func (b *Bot) Stop() {
	close(b.tasks)
	b.wg.Wait()
}
