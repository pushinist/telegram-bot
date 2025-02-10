package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
	"time"
)

func sendGifResponse(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, gifPath string) {
	if _, err := os.Stat(gifPath); os.IsNotExist(err) {
		log.Printf("Gif file not found: %s", gifPath)
		return
	}

	gif := tgbotapi.NewVideo(msg.Chat.ID, tgbotapi.FilePath(gifPath))
	gif.ReplyToMessageID = msg.MessageID

	maxRetries := 3
	for retry := 0; retry < maxRetries; retry++ {
		_, err := bot.Send(gif)
		if err != nil {
			if strings.Contains(err.Error(), "Too Many Requests") {
				retryAfter := 5 * time.Second
				if retry < maxRetries-1 {
					log.Printf("Rate limited. Waiting %v before retry %d/%d",
						retryAfter, retry+1, maxRetries)
					time.Sleep(retryAfter)
					continue
				}
			}
		}
		break
	}

}
