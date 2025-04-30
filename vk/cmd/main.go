package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vk/internal/clients/blogator"
	"vk/internal/clients/db"
	"vk/internal/clients/rabbitmq"
	vkClient "vk/internal/clients/vk"
	"vk/internal/config"
	event_consumer "vk/internal/consumer/event-consumer"
	"vk/internal/events/vk"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad()

	eventProccessor := vk.New(
		vkClient.New(cfg.VkToken),
		blogator.New(cfg.AtorToken),
		db.New(cfg.DbHost, cfg.DbPort),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	rmq, err := rabbitmq.New(cfg.RabbitUrl, cfg.RabbitQueue)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ client: %v", err)
	}
	defer rmq.Close()

	err = rmq.Consume(ctx, func(post rabbitmq.Response) error {
		return eventProccessor.SendPostToSubscribers(ctx, post)
	})

	if err != nil {
		log.Fatalf("Failed to set up RabbitMQ consumer: %v", err)
	}

	log.Print("service started")

	consumer := event_consumer.New(eventProccessor, eventProccessor)

	go func() {
		if err := consumer.Start(); err != nil {
			log.Printf("Telegram consumer stopped: %v", err)
			cancel()
		}
	}()

	log.Print("service started")

	// Ожидаем сигнал завершения
	sig := <-sigChan
	log.Printf("Received signal: %v, initiating graceful shutdown...", sig)
	cancel() // Отменяем контекст

	// Даем время на graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	log.Println("Closing RabbitMQ connection...")
	if err := rmq.Close(); err != nil {
		log.Printf("Error closing RabbitMQ: %v", err)
	}

	select {
	case <-shutdownCtx.Done():
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Println("Shutdown timed out, forcing exit")
		} else {
			log.Println("Shutdown complete")
		}
	}
}
