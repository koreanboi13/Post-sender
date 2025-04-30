package vk

import (
	"context"
	"fmt"
	"log"
	"vk/internal/clients/rabbitmq"
)

const (
	messangerType = "Telegram"
)

func (p *Processor) SendPostToSubscribers(ctx context.Context, post rabbitmq.Response) error {
	chatIDs, err := p.db.AllUsers(ctx, messangerType)
	if err != nil {
		return fmt.Errorf("failed to get subscribers: %w", err)
	}
	log.Println("TEST TEST TEST")
	if len(chatIDs) == 0 {
		log.Println("No subscribers to send post to")
		return nil
	}
	log.Println("Sending post to subscribers", post.Data)
	for _, ps := range post.Data {
		text := fmt.Sprintf(`
🗣 %s

📰 %s

🔗 %s
`, ps.Comment, ps.Title, ps.Link)
		log.Println("chatIDs: ", chatIDs)
		log.Println("post: ", text)
		for _, chatID := range chatIDs {
			log.Printf("Sending post to chat %d: %s", chatID, ps.Title)

			if err := p.vk.SendMessage(ctx, chatID, text); err != nil {
				log.Printf("Error sending message to chat %d: %v", chatID, err)
				continue
			}
		}
	}
	return nil
}
