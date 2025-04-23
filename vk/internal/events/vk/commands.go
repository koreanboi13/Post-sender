package vk

import (
	"context"
	"log"
	"strings"
)

const (
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%d", text, chatID)

	switch text {
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	default:
		return p.vk.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.vk.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.vk.SendMessage(ctx, chatID, msgHello)
}
