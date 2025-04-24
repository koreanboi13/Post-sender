package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	VkToken   string
	AtorToken string
	DbHost    string
	DbPort    string
}

func MustLoad() *config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env file load error")
	}

	cfg := &config{
		VkToken:   os.Getenv("VK_TOKEN"),
		AtorToken: os.Getenv("BLOGATOR_TOKEN"),
		DbHost:    os.Getenv("DB_HOST"),
		DbPort:    os.Getenv("DB_PORT"),
	}
	return cfg

}
