package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"
)

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

	//maxRetries := 3
	//for retry := 0; retry < maxRetries; retry++ {
	//	_, err := bot.Send(gif)
	//	slog.Info("Message sent")
	//	if err != nil {
	//		if strings.Contains(err.Error(), "Too Many Requests") {
	//			retryAfter := 5 * time.Second
	//			if retry < maxRetries-1 {
	//				log.Printf("Rate limited. Waiting %v before retry %d/%d",
	//					retryAfter, retry+1, maxRetries)
	//				time.Sleep(retryAfter)
	//				continue
	//			}
	//		}
	//	}
	//	break
	//}

}
