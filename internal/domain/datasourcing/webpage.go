package source

import (
	"time"

	"github.com/gocolly/colly/v2"
)

var USER_AGENTS = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
}

type WebPage struct {
	Name          string
	URL           string
	LastVisitedAt time.Time
	collector     *colly.Collector
	result        chan any
}

func NewWebPage(name string, URL string, collector *colly.Collector) *WebPage {
	return &WebPage{
		Name:      name,
		URL:       URL,
		collector: collector,
		result:    make(chan any),
	}
}

func (page *WebPage) Parse() <-chan any {
	page.collector.Visit(page.URL)

	go func() {
		page.collector.Wait()
		close(page.result)
	}()

	return page.result
}
