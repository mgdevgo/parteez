package scraping

import (
	"context"
	"log/slog"
	"sync"
)

type ScrapingService interface {
	Scrape(ctx context.Context) <-chan Result
}

type SourceService interface {
	ListSources() error
	Add(ctx context.Context) error
	Update(ctx context.Context) error
	Disable(ctx context.Context) error
}

type EventScrapingService struct {
	sources []Source
	logger  *slog.Logger
}

func NewEventScrapingService(sources []Source, logger *slog.Logger) *EventScrapingService {
	return &EventScrapingService{
		sources: sources,
		logger:  logger,
	}
}

func (service *EventScrapingService) Scrape(ctx context.Context) <-chan Result {
	resultChan := make(chan Result)
	var wg sync.WaitGroup
	wg.Add(len(service.sources))

	for _, source := range service.sources {
		inputChan, err := source.Parse(ctx)
		if err != nil {
			service.logger.Error("Failed to parse source", "error", err)
			continue
		}

		go func() {
			defer wg.Done()
			for item := range inputChan {
				resultChan <- item
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}
