package bot

import (
	"os"
)

func (b *Bot) startWorker(sigChan chan<- os.Signal) {
	defer b.wg.Done()

	for task := range b.tasks {
		if b.handlers.TextHandler.CanHandle(task.Message) {
			err := b.handlers.TextHandler.Handle(task)
			if err != nil {
				sigChan <- os.Interrupt
			}
			break
		}

		if b.handlers.GifHandler.CanHandle(task.Message) {
			err := b.handlers.GifHandler.Handle(task)
			if err != nil {
				sigChan <- os.Interrupt
			}
			break
		}
	}
}
