package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	slog.SetDefault(Logger)
}
