package main

import (
	"api/internal/clients/blogator"
	"api/internal/clients/db"
	tgClient "api/internal/clients/telegram"
	"api/internal/config"
	event_consumer "api/internal/consumer/event-consumer"
	"api/internal/events/telegram"
	"log"
)

func main() {
	cfg := config.MustLoad()

	eventProccessor := telegram.New(
		tgClient.New(cfg.TgHost, cfg.TgToken),
		blogator.New(cfg.AtorToken),
		db.New(cfg.DbHost, cfg.DbPort),
	)

	log.Print("service started")

	consumer := event_consumer.New(eventProccessor, eventProccessor, cfg.BatchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service stopped", err)
	}
}
