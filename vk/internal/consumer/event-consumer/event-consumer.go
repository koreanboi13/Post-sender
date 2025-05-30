package event_consumer

import (
	"context"
	"log"
	"time"
	"vk/internal/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
}

func New(fetcher events.Fetcher, processor events.Processor) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(context.Background())
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(context.Background(), gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(ctx context.Context, events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(ctx, event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
