package web

import (
	"fmt"
	"iditusi/internal/parsers/result"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/v2"
)

var RuporCollector *colly.Collector

func init() {
	RuporCollector = colly.NewCollector(
		colly.UserAgent(USER_AGENT),
		colly.Async(true),
		// colly.Debugger(&debug.LogDebugger{}),
		// colly.AllowedDomains(Domain),

	)
}

func NewRuporEventsCrWebPage() *CrawlableWebPage {
	const name = "ptichka"
	const URL = "https://listing.events/spb"

	collyCollector :=

		collyCollector.Limit(&colly.LimitRule{
			RandomDelay: 2 * time.Second,
		})

	collyCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL: ", r.URL.String())
	})

	collyCollector.OnHTML("article > div.event-information-banner__afisha > a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		collyCollector.Visit(link)
	})

	collyCollector.OnHTML("article[itemscope]", func(e *colly.HTMLElement) {
		var event result.Data

		// Age restriction
		age := e.ChildText("h1.event-information__name > span.age-limit")
		event.AgeRestriction, _ = strconv.Atoi(strings.Trim(age, "+"))

		// Name
		event.Tittle = strings.TrimSpace(strings.TrimRight(e.ChildText("h1.event-information__name"), age))

		// Description
		event.Description = strings.TrimSpace(e.ChildText("div.event-information__description"))

		// ArtworkURL
		event.ArtworkURL = strings.Trim(e.ChildAttr("div.event-information-banner", "style"), "background-image: url()")

		// Genres
		event.Genres = make([]string, 0)
		e.ForEach("div.bottom-tags-wrapper > span.btn-tag-genres", func(i int, e *colly.HTMLElement) {
			event.Genres = append(event.Genres, strings.TrimSpace(e.Text))
		})

		// LineUp
		e.ForEach(`ul.program-block__list > li.program-block__item[itemprop="performer"]`, func(i int, e *colly.HTMLElement) {
			event.LineUp = append(event.LineUp, squashSpace(e.Text))
		})
		// artists := make([]models.LineUpArtist, 0)
		// e.ForEach(`ul.program-block__list > li.program-block__item[itemprop="performer"]`, func(i int, e *colly.HTMLElement) {
		// 	artists = append(artists, models.LineUpArtist{Name: squashSpace(e.Text)})
		// })
		// event.LineUp = models.LineUp{
		// 	{"main", artists},
		// }

		// Dates
		loc, _ := time.LoadLocation("Europe/Moscow")
		// Start date and time
		event.StartTime, _ = time.ParseInLocation(
			time.DateTime,
			fmt.Sprintf(
				"%s %s:00",
				e.ChildAttr(`p > time[itemprop="startDate"]`, "datetime"),
				strings.Fields(e.ChildText(`p > time[itemprop="startDate"]`))[3],
			),
			loc,
		)
		// End date. If not presented time.Time zero value
		event.EndTime, _ = time.ParseInLocation(
			time.DateTime,
			fmt.Sprintf(
				"%s %s:00",
				e.ChildAttr(`p > time[itemprop="endDate"]`, "datetime"),
				strings.TrimSpace(e.ChildText(`p > time[itemprop="endDate"]`)),
			),
			loc,
		)

		// Venue
		event.Location.Name = strings.TrimSpace(e.ChildText("p.venue-name > a"))

		address := strings.Split(e.ChildText("p.venue-adress"), "м.") // typo is left intentionally
		event.Location.Address = squashSpace(address[0])

		stations := make([]string, 0)
		if len(address) > 1 {
			metro := strings.Split(address[1], ",")
			for _, i := range metro {
				stations = append(stations, strings.TrimSpace(i))
			}
		}
		event.Location.MetroStations = stations

		// Tickets
		// TODO: get full url instead of shorten
		event.TicketsURL = e.ChildAttr("a.buy-btn", "href")

		output <- event
	})

	// Create another collector to scrape event details
	// details := p.collector.Clone()

	return NewCrawlableWebPage(name, URL, collector)
}
