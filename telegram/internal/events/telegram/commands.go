package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	SubscribeCmd   = "/subscribe"
	UnsubscribeCmd = "/unsubscribe"
	HelpCmd        = "/help"
	StartCmd       = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s' with chatID: %d", text, username, chatID)

	switch text {
	case SubscribeCmd:
		return p.subscribe(ctx, chatID)
	case UnsubscribeCmd:
		return p.unsubscribe(ctx, chatID)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) subscribe(ctx context.Context, chatID int) error {
	res, err := p.db.ChatExists(ctx, chatID)
	if err != nil {
		return fmt.Errorf("can't check if chat exists: %w", err)
	}
	if res {
		return p.tg.SendMessage(ctx, chatID, msgAlreadySubscribed)
	}

	if err = p.db.SaveUser(ctx, chatID); err != nil {
		return fmt.Errorf("can't save user: %w", err)
	}
	return p.tg.SendMessage(ctx, chatID, msgSubscribed)
}

func (p *Processor) unsubscribe(ctx context.Context, chatID int) error {
	exists, err := p.db.ChatExists(ctx, chatID)
	if err != nil {
		return fmt.Errorf("can't check if chat exists: %w", err)
	}

	if !exists {
		return p.tg.SendMessage(ctx, chatID, "You are already unsubscribed!")
	}

	if err := p.db.DeleteUser(ctx, chatID); err != nil {
		return fmt.Errorf("can't unsubscribe: %w", err)
	}

	return p.tg.SendMessage(ctx, chatID, msgUnsubscribedSuccess)
}

func (p *Processor) Post(ctx context.Context, chatID int) (err error) {
	twoDaysAgo := time.Now().UTC().AddDate(0, 0, -2)
	log.Printf("Fetching posts since %v", twoDaysAgo)

	items, err := p.ator.GetNewItems(ctx, twoDaysAgo)
	if err != nil {
		return fmt.Errorf("can't get new posts: %w", err)
	}

	if len(items) == 0 {
		return p.tg.SendMessage(ctx, chatID, msgNoPosts)
	}

	log.Printf("Got %d new posts", len(items))

	for _, item := range items {
		text := fmt.Sprintf(`
📰 %s

🔗 %s
`, item.Title, item.Link)

		log.Printf("Sending message to chat %d: %s", chatID, item.Title)
		if err := p.tg.SendMessage(ctx, chatID, text); err != nil {
			log.Printf("Error sending message: %v", err)
			return fmt.Errorf("can't send message: %w", err)
		}
	}
	return nil
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}
