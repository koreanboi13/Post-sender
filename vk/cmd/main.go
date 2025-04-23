package main

import (
	"log"
	"vk/internal/clients/blogator"
	vkClient "vk/internal/clients/vk"
	"vk/internal/config"
	event_consumer "vk/internal/consumer/event-consumer"
	"vk/internal/events/vk"
)

func main() {
	cfg := config.MustLoad()
	eventProccessor := vk.New(vkClient.New(cfg.VkToken), blogator.New(cfg.AtorToken))

	log.Print("service started")

	consumer := event_consumer.New(eventProccessor, eventProccessor)

	if err := consumer.Start(); err != nil {
		log.Fatal("service stopped", err)
	}

}
