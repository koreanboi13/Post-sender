package vk

import (
	"context"
	"fmt"
	"log"
	"strings"
)

const (
	SubscribeCmd   = "/subscribe"
	UnsubscribeCmd = "/unsubscribe"
	HelpCmd        = "/help"
	StartCmd       = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%d", text, chatID)

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
		return p.vk.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) subscribe(ctx context.Context, chatID int) error {
	res, err := p.db.ChatExists(ctx, chatID)
	if err != nil {
		return fmt.Errorf("can't check if chat exists: %w", err)
	}
	if res {
		return p.vk.SendMessage(ctx, chatID, msgAlreadySubscribed)
	}

	if err = p.db.SaveUser(ctx, chatID); err != nil {
		return fmt.Errorf("can't save user: %w", err)
	}
	return p.vk.SendMessage(ctx, chatID, msgSubscribed)
}

func (p *Processor) unsubscribe(ctx context.Context, chatID int) error {
	exists, err := p.db.ChatExists(ctx, chatID)
	if err != nil {
		return fmt.Errorf("can't check if chat exists: %w", err)
	}

	if !exists {
		return p.vk.SendMessage(ctx, chatID, "You are already unsubscribed!")
	}

	if err := p.db.DeleteUser(ctx, chatID); err != nil {
		return fmt.Errorf("can't unsubscribe: %w", err)
	}

	return p.vk.SendMessage(ctx, chatID, msgUnsubscribedSuccess)
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.vk.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.vk.SendMessage(ctx, chatID, msgHello)
}
