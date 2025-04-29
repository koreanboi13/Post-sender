package main

import (
	"api/internal/clients/blogator"
	"api/internal/clients/db"
	"api/internal/clients/rabbitmq"
	tgClient "api/internal/clients/telegram"
	"api/internal/config"
	event_consumer "api/internal/consumer/event-consumer"
	"api/internal/events/telegram"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настройка обработки сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	cfg := config.MustLoad()

	eventProccessor := telegram.New(
		tgClient.New(cfg.TgHost, cfg.TgToken),
		blogator.New(cfg.AtorToken),
		db.New(cfg.DbHost, cfg.DbPort),
	)

	// Инициализируем RabbitMQ клиент
	rmq, err := rabbitmq.New(cfg.RabbitUrl, cfg.RabbitQueue)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ client: %v", err)
	}
	defer rmq.Close()

	// Настраиваем потребителя RabbitMQ
	err = rmq.Consume(ctx, func(post rabbitmq.Response) error {
		return eventProccessor.SendPostToSubscribers(ctx, post)
	})
	if err != nil {
		log.Fatalf("Failed to set up RabbitMQ consumer: %v", err)
	}

	// Создаем consumer для Telegram
	consumer := event_consumer.New(eventProccessor, eventProccessor, cfg.BatchSize)

	// Запускаем Telegram consumer в отдельной goroutine
	go func() {
		if err := consumer.Start(); err != nil {
			log.Printf("Telegram consumer stopped: %v", err)
			cancel() // Отменяем контекст для остановки всех компонентов
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

	// Ждем завершения или тайм-аута
	select {
	case <-shutdownCtx.Done():
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Println("Shutdown timed out, forcing exit")
		} else {
			log.Println("Shutdown complete")
		}
	}
}
