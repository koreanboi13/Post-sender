package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	queue    string
	stopChan chan struct{}  // Канал для сигнала остановки
	wg       sync.WaitGroup // Группа ожидания для горутин
}

func New(Url string, queue string) (*Client, error) {
	conn, err := amqp.Dial(Url)
	if err != nil {
		return nil, fmt.Errorf("can`t connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("can`t create channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("can`t declare queue: %w", err)
	}
	return &Client{
		conn:     conn,
		ch:       ch,
		queue:    q.Name,
		stopChan: make(chan struct{}),
	}, nil
}

func (c *Client) Consume(ctx context.Context, handler func(post Response) error) error {
	msgs, err := c.ch.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}
	log.Println("Consumer registered")
	log.Println("new messages", msgs)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Println("Context canceled, stopping RabbitMQ consumer")
				return
			case <-c.stopChan:
				log.Println("Received stop signal, stopping RabbitMQ consumer")
				return
			case d, ok := <-msgs:
				if !ok {
					log.Println("RabbitMQ channel closed")
					return
				}

				var post Response
				if err := json.Unmarshal(d.Body, &post); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					d.Nack(false, false) // Reject message
					continue
				}

				if err := handler(post); err != nil {
					log.Printf("Error handling message: %v", err)
					d.Nack(false, true) // Requeue message
					continue
				}

				d.Ack(false) // Acknowledge message
			}
		}
	}()

	return nil
}

// Stop останавливает потребителя и ждет завершения всех горутин
func (c *Client) Stop() {
	close(c.stopChan)
	c.wg.Wait() // Ждем завершения всех горутин
}

// Close закрывает соединение с RabbitMQ
func (c *Client) Close() error {
	c.Stop() // Сначала останавливаем потребителя

	if err := c.ch.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}
