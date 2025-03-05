package scraping

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Rule func(c *colly.Collector, result chan<- any)

var scrapingRules = map[string]Rule{
	"listing.events": func(c *colly.Collector, result chan<- any) {
		c.OnHTML("article > div.event-information-banner__afisha > a", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			c.Visit(e.Request.AbsoluteURL(link))
		})
		c.OnHTML("article[itemscope]", func(e *colly.HTMLElement) {
			data := make(map[string]any)
			age := e.ChildText("h1.event-information__name > span.age-limit")
			data["ageRestriction"], _ = strconv.Atoi(strings.Trim(age, "+"))
			data["tittle"] = strings.TrimSpace(strings.TrimRight(e.ChildText("h1.event-information__name"), age))
			data["description"] = strings.TrimSpace(e.ChildText("div.event-information__description"))
			data["artworkURL"] = strings.Trim(e.ChildAttr("div.event-information-banner", "style"), "background-image: url()")
			genres := make([]string, 0)
			e.ForEach("div.bottom-tags-wrapper > span.btn-tag-genres", func(i int, e *colly.HTMLElement) {
				genres = append(genres, strings.TrimSpace(e.Text))
			})
			data["genres"] = genres
			lineup := make([]string, 0)
			e.ForEach(`ul.program-block__list > li.program-block__item[itemprop="performer"]`, func(i int, e *colly.HTMLElement) {
				lineup = append(lineup, squashSpace(e.Text))
			})
			data["lineup"] = lineup

			loc, _ := time.LoadLocation("Europe/Moscow")
			data["startTime"], _ = time.ParseInLocation(
				time.DateTime,
				fmt.Sprintf(
					"%s %s:00",
					e.ChildAttr(`p > time[itemprop="startDate"]`, "datetime"),
					strings.Fields(e.ChildText(`p > time[itemprop="startDate"]`))[3],
				),
				loc,
			)
			data["endTime"], _ = time.ParseInLocation(
				time.DateTime,
				fmt.Sprintf(
					"%s %s:00",
					e.ChildAttr(`p > time[itemprop="endDate"]`, "datetime"),
					strings.TrimSpace(e.ChildText(`p > time[itemprop="endDate"]`)),
				),
				loc,
			)
			location := make(map[string]any)
			location["name"] = strings.TrimSpace(e.ChildText("p.venue-name > a"))

			address := strings.Split(e.ChildText("p.venue-adress"), "м.") // typo is left intentionally
			location["address"] = squashSpace(address[0])

			stations := make([]string, 0)
			if len(address) > 1 {
				metro := strings.Split(address[1], ",")
				for _, i := range metro {
					stations = append(stations, strings.TrimSpace(i))
				}
			}
			location["metroStations"] = stations

			data["ticketsURL"] = e.ChildAttr("a.buy-btn", "href")

			result <- data
		})
	},
	"blunk.ru": func(c *colly.Collector, result chan<- any) {},
}

// squashSpace turn many spaces into one
func squashSpace(text string) string {
	return strings.TrimSpace(regexp.MustCompile(`\s{2,}`).ReplaceAllString(text, " "))
}
