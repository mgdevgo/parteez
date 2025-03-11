package scraping

import (
	"context"
	"log"
	"log/slog"

	"parteez/internal/domain/events"
)

type Agent struct {
	scrapingService  ScrapingService
	eventCrudService events.EventCrudService
	logger           *slog.Logger
}

func NewAgent(scrapingService ScrapingService, eventCrudService events.EventCrudService) *Agent {
	return &Agent{
		scrapingService:  scrapingService,
		eventCrudService: eventCrudService,
	}
}

func (a *Agent) Run(ctx context.Context) {
	results := a.scrapingService.Scrape(ctx)

	for result := range results {
		event, err := a.eventCrudService.CreateDraft(ctx)
		if err != nil {
			log.Printf("error creating draft: %v", err)
			continue
		}
		_ = result
		a.eventCrudService.Update(ctx, event.ID, events.EventUpdate{})
	}
}

func (a *Agent) Stop() {

}
