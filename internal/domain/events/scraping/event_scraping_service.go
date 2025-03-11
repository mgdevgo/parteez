package scraping

import (
	"context"
	"log/slog"
	"sync"
)

var USER_AGENTS = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
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
		inputChan := source.Parse()
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
