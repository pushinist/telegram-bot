package main

import (
	"fmt"
	"github.com/pushinist/telegram-bot/internal/bot"
	"github.com/pushinist/telegram-bot/internal/config"
	"github.com/pushinist/telegram-bot/pkg/logger"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Starting Telegram Bot")
	logger.Init()
	slog.Info("Logger initialized")
	cfg, err := config.Load()
	if err != nil {
		slog.Error(fmt.Sprintf("Error loading config: %v", err))
	}
	slog.Info("Config loaded")

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
