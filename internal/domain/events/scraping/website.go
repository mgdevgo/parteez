package scraping

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gocolly/colly/v2"
)

type Website struct {
	domain    string
	url       string
	collector *colly.Collector
	result    chan any
	debug     bool
}

func NewWebsite(URL string, collector *colly.Collector) (*Website, error) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	name := parsedURL.Host

	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent(USER_AGENTS[0]),
	)

	if err := c.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second,
	}); err != nil {
		return nil, err
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("%s | visiting URL: %s\n", name, r.URL.String())
	})

	ruleFunc, ok := scrapingRules[name]
	if !ok {
		return nil, fmt.Errorf("source %s not found", name)
	}

	site := Website{
		domain:    name,
		url:       URL,
		collector: collector,
		result:    make(chan any),
	}

	ruleFunc(c, site.result)

	return &site, nil
}

func (page *Website) Parse() chan any {
	page.collector.Visit(page.url)

	go func() {
		page.collector.Wait()
		close(page.result)
	}()

	return page.result
}

func (site *Website) ID() string {
	return site.domain
}

func (site *Website) Debug(debug bool) {
	site.debug = debug
}
