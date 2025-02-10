package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

	taskChan := make(chan MessageTask, 100)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go messageWorker(taskChan, &wg)
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

func messageWorker(taskChan <-chan MessageTask, wg *sync.WaitGroup) {
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
					sendGifResponse(task.Bot, task.Message, gifPath)
					break
				}
			}
		}
		if task.Message.Animation != nil {
			log.Printf("Recieved GIF with the unique ID: %s", task.Message.Animation.FileUniqueID)

			if responsePath, exists := gifTriggers[task.Message.Animation.FileUniqueID]; exists {
				sendGifResponse(task.Bot, task.Message, responsePath)
			}
		}

		messageText := task.Message.Text
		chatID := task.Message.Chat.ID

		for phrase, gifPath := range triggerGifs {
			if containsTrigger(messageText, phrase) {
				gif := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(gifPath))
				gif.ReplyToMessageID = task.Message.MessageID
				_, err := task.Bot.Send(gif)
				if err != nil {
					log.Println(err)
				}
				break
			}
		}
	}
}

func sendGifResponse(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, gifPath string) {
	if _, err := os.Stat(gifPath); os.IsNotExist(err) {
		log.Printf("Gif file not found: %s", gifPath)
		return
	}

	gif := tgbotapi.NewDocument(msg.Chat.ID, tgbotapi.FilePath(gifPath))
	gif.ReplyToMessageID = msg.MessageID
	_, err := bot.Send(gif)
	if err != nil {
		log.Println(err)
	}

}

func containsTrigger(message, trigger string) bool {
	return strings.Contains(strings.ToLower(message), strings.ToLower(trigger))
}
