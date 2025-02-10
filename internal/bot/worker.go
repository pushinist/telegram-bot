package bot

import (
	"fmt"
	"log/slog"
)

func (b *Bot) startWorker() {
	defer b.wg.Done()

	for task := range b.tasks {
		slog.Info(fmt.Sprintf("task message: %v", task.Message))
		if b.handlers.TextHandler.CanHandle(task.Message) {
			slog.Info(fmt.Sprintf("Message read: %v", task.Message))
			err := b.handlers.TextHandler.Handle(task)
			if err != nil {
				slog.Error(err.Error())
			}
			break
		}

		if b.handlers.GifHandler.CanHandle(task.Message) {
			err := b.handlers.GifHandler.Handle(task)
			slog.Info(fmt.Sprintf("Message read: %v", task.Message))
			if err != nil {
				//sigChan <- os.Interrupt
			}
			break
		}
	}
}
