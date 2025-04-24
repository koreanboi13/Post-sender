package vk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"vk/internal/clients/blogator"
	"vk/internal/clients/db"
	"vk/internal/clients/vk"
	"vk/internal/events"
)

type Processor struct {
	vk   *vk.Client
	ator *blogator.Client
	db   *db.Client
}

type Meta struct {
	PeerID int
	FromID int
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *vk.Client, blog *blogator.Client, db *db.Client) *Processor {
	return &Processor{
		vk:   client,
		ator: blog,
		db:   db,
	}
}

func (p *Processor) Fetch(ctx context.Context) ([]events.Event, error) {
	vkUpdates, err := p.vk.Updates(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get updates: %w", err)
	}

	if len(vkUpdates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(vkUpdates))

	for _, u := range vkUpdates {
		event, err := processUpdate(u)
		if err != nil {
			log.Printf("Error processing update: %v", err)
			continue
		}
		res = append(res, event)
	}

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	default:
		return fmt.Errorf("can't process event")
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := getMeta(event)
	if err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	if err := p.doCmd(ctx, event.Text, meta.PeerID); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func getMeta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("can't get Meta")
	}

	return res, nil
}

func processUpdate(upd vk.LongPollUpdate) (events.Event, error) {
	if upd.Type != "message_new" {
		return events.Event{}, fmt.Errorf("not a message event: %s", upd.Type)
	}

	text := upd.Object.Text
	fromID := upd.Object.FromId

	return events.Event{
		Type: events.Message,
		Text: text,
		Meta: Meta{
			PeerID: fromID,
			FromID: fromID,
		},
	}, nil
}
