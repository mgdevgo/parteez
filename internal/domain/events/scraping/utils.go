package scraping

import (
	"log/slog"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

var USER_AGENTS = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
}

func commonCollyOptions(logger *slog.Logger) []colly.CollectorOption {
	return []colly.CollectorOption{
		colly.Async(true),
		// colly.UserAgent(USER_AGENTS[0]),
		logRequest(logger),
		extensions.RandomUserAgent,
	}
}

func logRequest(logger *slog.Logger) colly.CollectorOption {
	return func(c *colly.Collector) {
		c.OnRequest(func(r *colly.Request) {
			logger.Info("Requesting website", "url", r.URL.String())
		})
	}
}

// squashSpace turn many spaces into one
func squashSpace(text string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(text, " "))
}
