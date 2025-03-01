package source

import (
	"sync"

	"github.com/gocolly/colly/v2"
)

type CrowdSourceService struct {
	sources                 []Source
	maxRetries              int
	webPageCllectos         []*colly.Collector
	telegramChannelsParsers []any
	errorHandler            func(error)
	output                  chan any
	duplicates              map[string]struct{}
	duplicatesLock          sync.RWMutex
}

func (service *CrowdSourceService) AddSource() error {
	panic("not implemented")
}

func (service *CrowdSourceService) ListSources() error {
	panic("not implemented")
}

func (service *CrowdSourceService) Parse() error {
	panic("not implemented")
}

func (service *CrowdSourceService) FetchData() chan any {
	wg := sync.WaitGroup{}

	for _, source := range service.sources {
		wg.Add(1)
		output := source.Parse()
		go func() {
			defer wg.Done()
			for item := range output {
				service.duplicatesLock.Lock()

				key := source.ID()
				_, ok := service.duplicates[key]
				if !ok {
					service.duplicates[key] = struct{}{}
					service.output <- item
				}

				service.duplicatesLock.Unlock()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(service.output)
	}()

	return service.output
}
