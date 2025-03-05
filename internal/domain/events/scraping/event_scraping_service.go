package scraping

import (
	"context"
	"sync"
)

var USER_AGENTS = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
}

type EventScrapingService struct {
	sources []Source
	output  chan any
}

func NewEventScrapingService(sources []Source) *EventScrapingService {
	return &EventScrapingService{
		sources: sources,
		output:  make(chan any),
	}
}

func (service *EventScrapingService) Scrape(ctx context.Context) chan any {
	wg := sync.WaitGroup{}

	for _, source := range service.sources {
		wg.Add(1)
		output := source.Parse()
		go func() {
			defer wg.Done()
			for item := range output {
				select {
				case <-ctx.Done():
					return
				case service.output <- item:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(service.output)
	}()

	return service.output
}
