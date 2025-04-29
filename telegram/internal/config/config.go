package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	TgToken     string
	TgHost      string
	AtorToken   string
	BatchSize   int
	DbHost      string
	DbPort      string
	RabbitUrl   string
	RabbitQueue string
}

func MustLoad() *config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env file load error")
	}
	size, err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	if err != nil {
		log.Fatal("can`t get batch size")
	}

	cfg := &config{
		TgToken:     os.Getenv("TELEGRAM_TOKEN"),
		TgHost:      os.Getenv("TELEGRAM_HOST"),
		AtorToken:   os.Getenv("BLOGATOR_TOKEN"),
		DbHost:      os.Getenv("DB_HOST"),
		DbPort:      os.Getenv("DB_PORT"),
		BatchSize:   size,
		RabbitUrl:   os.Getenv("RABBITMQ_URL"),
		RabbitQueue: os.Getenv("RABBITMQ_QUEUE"),
	}
	return cfg

}
