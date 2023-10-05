package event_consumer

import (
	"log/slog"
	"read-adviser-bot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
	logger    *slog.Logger
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int, logger *slog.Logger) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		logger:    logger,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			c.logger.Error("[ERR] consumer: ", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(&gotEvents); err != nil {
			c.logger.Error("can't handle events: ", slog.String("error", err.Error()))

			continue
		}
	}
}

func (c *Consumer) handleEvents(events *[]events.Event) error {
	for _, event := range *events {
		// TODO: retry/backup, error counter
		if err := c.processor.Process(event); err != nil {
			c.logger.Error("can't handle event: ", slog.String("error", err.Error()))

			continue
		}
	}

	return nil
}
