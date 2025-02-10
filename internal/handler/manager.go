package handler

type Manager struct {
	GifHandler  *GifHandler
	TextHandler *TextHandler
}

func NewManager() *Manager {
	return &Manager{
		GifHandler:  NewGifHandler(),
		TextHandler: NewTextHandler(),
	}
}
