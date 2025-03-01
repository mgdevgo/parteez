package collectors

import "github.com/gocolly/colly/v2"

func Blank(URL string) *colly.Collector {
	// name := "blank"
	// URL := "https://blankclub.ru"

	collector := colly.NewCollector(
	// TODO: настроить скрапер
	)

	return collector
}
