package scraping

import "context"

type Data map[string]any

type ScrapingService interface {
	// FetchData(sources []int) chan any
	Scrape(ctx context.Context) chan Data
}

// type EventSourceService interface {
// 	ListSources() error
// 	Add(ctx context.Context, source Source) error
// 	Update(ctx context.Context, source Source) error
// 	Disable(ctx context.Context, sourceId int) error
// }
