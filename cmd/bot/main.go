package main

import (
	"fmt"
	"github.com/pushinist/telegram-bot/internal/bot"
	"github.com/pushinist/telegram-bot/internal/config"
	"github.com/pushinist/telegram-bot/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		slog.Error(fmt.Sprintf("Error loading config: %v", err))
	}

	tgBot, err := bot.New(cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("Error creating telegram bot: %v", err))
	}
	slog.Info("Bot started")
	go tgBot.Start()
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	slog.Info("Shutting down")

	tgBot.Stop()
}
