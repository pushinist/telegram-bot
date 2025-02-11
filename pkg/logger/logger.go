package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init() {
	//logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	slog.Error("failed to open log file", "error", err)
	//	return
	//}

	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(Logger)
}
