package main

import (
	"github.com/pushinist/telegram-bot/internal/bot"
	"github.com/pushinist/telegram-bot/internal/config"
	"github.com/pushinist/telegram-bot/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Init()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	bot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	go bot.Start(sigChan)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	bot.Stop()
}
