package bot

import (
	"fmt"
	"log/slog"

	"github.com/pushinist/telegram-bot/internal/handler"
)

func (b *Bot) startWorker(id int) {
	defer b.wg.Done()

	slog.Info(fmt.Sprintf("worker %d starting", id))

	for task := range b.tasks {
		slog.Info(fmt.Sprintf("worker %d processing message: %v", id, handler.ParseMessage(task.Message)))

		if err := b.handlers.Handle(task); err != nil {
			slog.Error(err.Error())
		}

		//if b.handlers.TextHandler.CanHandle(task.Message) {
		//	slog.Info(fmt.Sprintf("worker %d processing text message: %v", task.Message.Text))
		//	if err := b.handlers.TextHandler.Handle(task); err != nil {
		//		slog.Error(err.Error())
		//	}
		//}
		//
		//if b.handlers.GifHandler.CanHandle(task.Message) {
		//	slog.Info(fmt.Sprintf("worker %d processing gif message: %v", task.Message.Animation.FileUniqueID))
		//	if err := b.handlers.GifHandler.Handle(task); err != nil {
		//		slog.Error(err.Error())
		//	}
		//
		//}
	}
}
