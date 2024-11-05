package web

import "github.com/gocolly/colly/v2"

func NewBlankPage() *CrawlableWebPage {
	name := "blank"
	URL := "https://blankclub.ru"

	collector := colly.NewCollector(
	// TODO: настроить скрапер
	)

	return NewCrawlableWebPage(
		name,
		URL,
		collector,
	)
}
