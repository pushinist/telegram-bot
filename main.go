package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RateLimiter struct {
	lastRequest time.Time
	minInterval time.Duration
}

func newRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{
		lastRequest: time.Now(),
		minInterval: interval,
	}
}

func (rl RateLimiter) Wait() {
	elapsed := time.Since(rl.lastRequest)
	if elapsed < rl.minInterval {
		time.Sleep(rl.minInterval - elapsed)
	}
	rl.lastRequest = time.Now()
}

type MessageTask struct {
	Message *tgbotapi.Message
	Bot     *tgbotapi.BotAPI
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		fmt.Println("TELGERAM_BOT_TOKEN environment variable not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	rateLimiter := newRateLimiter(time.Second)

	taskChan := make(chan MessageTask, 300)

	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go messageWorker(taskChan, &wg, rateLimiter)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Chat.Type != "private" &&
			update.Message.Chat.Type != "group" &&
			update.Message.Chat.Type != "supergroup" {
			log.Printf("Skipping message from unsupported chat type: %s",
				update.Message.Chat.Type)
			continue
		}

		taskChan <- MessageTask{
			Message: update.Message,
			Bot:     bot,
		}
	}
}

func messageWorker(taskChan <-chan MessageTask, wg *sync.WaitGroup, rl *RateLimiter) {
	defer wg.Done()

	triggerGifs := map[string]string{
		"разрабы дауны": "./gifs/devs.gif",
		"разраб даун":   "./gifs/devs.gif",
		"даун":          "./gifs/devs.gif",
	}

	gifTriggers := map[string]string{
		"AgADzQIAAmRjXFM": "./gifs/devs.gif",
	}

	for task := range taskChan {

		if task.Message.Text != "" {
			messageText := task.Message.Text
			for phrase, gifPath := range triggerGifs {
				if containsTrigger(messageText, phrase) {
					sendGifResponse(task.Bot, task.Message, rl, gifPath)
					break
				}
			}
		}
		if task.Message.Animation != nil {
			log.Printf("Recieved GIF with the unique ID: %s", task.Message.Animation.FileUniqueID)

			if responsePath, exists := gifTriggers[task.Message.Animation.FileUniqueID]; exists {
				sendGifResponse(task.Bot, task.Message, rl, responsePath)
				break
			}
		}
	}
}

func sendGifResponse(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, rl *RateLimiter, gifPath string) {
	if _, err := os.Stat(gifPath); os.IsNotExist(err) {
		log.Printf("Gif file not found: %s", gifPath)
		return
	}

	rl.Wait()

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

func containsTrigger(message, trigger string) bool {
	return strings.Contains(strings.ToLower(message), strings.ToLower(trigger))
}
