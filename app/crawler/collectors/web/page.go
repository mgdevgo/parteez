package web

import (
	"iditusi/internal/parsers/result"

	"github.com/gocolly/colly/v2"
)

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0"

type CrawlableWebPage struct {
	name      string
	url       string
	collector *colly.Collector
	output    chan result.Data
}

func NewCrawlableWebPage(name string, URL string, collector *colly.Collector) *CrawlableWebPage {
	return &CrawlableWebPage{
		name:      name,
		url:       URL,
		collector: collector,
		output:    make(chan result.Data),
	}
}

func (p *CrawlableWebPage) Parse() <-chan result.Data {
	p.collector.Visit(p.url)

	go func() {
		p.collector.Wait()
		close(p.output)
	}()

	return p.output
}
