package config

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	Env      string
	BotToken string
	Workers  int
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	if _, exists := os.LookupEnv("TELEGRAM_BOT_TOKEN"); !exists {
		return nil, errors.New("telegram bot token not set")
	}
	return &Config{
		Env:      getEnv("ENV", "dev"),
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		Workers:  getEnvAsInt("WORKERS", 3),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	stringValue := os.Getenv(key)
	integerValue, err := strconv.Atoi(stringValue)
	if err != nil {
		return defaultValue
	}
	return integerValue
}
