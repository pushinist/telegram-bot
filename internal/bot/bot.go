package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pushinist/telegram-bot/internal/config"
	"github.com/pushinist/telegram-bot/internal/handler"
	"github.com/pushinist/telegram-bot/internal/model"
	"os"
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

func (b *Bot) Start(sigChan chan<- os.Signal) {
	for i := 0; i < b.workers; i++ {
		b.wg.Add(1)
		go b.startWorker(sigChan)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		b.tasks <- model.MessageTask{
			Message: update.Message,
			Bot:     b.api,
		}
	}
}

func (b *Bot) Stop() {
	close(b.tasks)
	b.wg.Wait()
}
