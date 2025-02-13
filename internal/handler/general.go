package handler

import (
	"fmt"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ParseMessage(message *tgbotapi.Message) string {
	if message.Text != "" {
		return message.Text
	}
	if message.Animation != nil {
		return message.Animation.FileUniqueID
	}

	return fmt.Sprintf("%v", message)
}

func sendGifResponse(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, gifPath string) error {
	if _, err := os.Stat(gifPath); os.IsNotExist(err) {
		slog.Error(fmt.Sprintf("Gif file not found: %s", gifPath))
		return err
	}

	gif := tgbotapi.NewVideo(msg.Chat.ID, tgbotapi.FilePath(gifPath))
	gif.ReplyToMessageID = msg.MessageID
	_, err := bot.Send(gif)
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("Message sent")
	return nil

}
