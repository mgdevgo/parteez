package scraping

import (
	"log/slog"
	"net/url"

	"github.com/gocolly/colly/v2"
)

type Website struct {
	domain    string
	url       string
	collector *colly.Collector
	result    chan Result
	debug     bool
	logger    *slog.Logger
}

func NewWebsite(URL string, logger *slog.Logger) (*Website, error) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	name := parsedURL.Host

	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent(USER_AGENTS[0]),
	)

	// limitRule := &colly.LimitRule{
	// 	RandomDelay: 2 * time.Second,
	// }
	// if err := limitRule.Init(); err != nil {
	// 	return nil, err
	// }
	// if err := c.Limit(limitRule); err != nil {
	// 	return nil, err
	// }

	c.OnRequest(func(r *colly.Request) {
		logger.Info("Requesting website", "url", r.URL.String())
	})

	site := &Website{
		domain:    name,
		url:       URL,
		collector: c,
		result:    make(chan Result),
	}

	return site, nil
}

func (w *Website) Collector() *colly.Collector {
	return w.collector
}

func (w *Website) Result(result Result) {
	w.result <- result
}

func (w *Website) Parse() chan Result {
	w.collector.Visit(w.url)

	go func() {
		w.collector.Wait()
		close(w.result)
	}()

	return w.result
}

func (w *Website) ID() string {
	return w.domain
}

func (w *Website) Debug(debug bool) {
	w.debug = debug
}
