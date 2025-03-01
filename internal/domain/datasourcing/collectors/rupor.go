package collectors

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func Rupor(URL string) *colly.Collector {
	const name = "ptichka"
	// const URL = "https://listing.events/spb"

	c := colly.NewCollector(
		// colly.UserAgent(USER_AGENT),
		colly.Async(true),
		// colly.Debugger(&debug.LogDebugger{}),
		// colly.AllowedDomains(Domain),

	)

	c.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL: ", r.URL.String())
	})

	c.OnHTML("article > div.event-information-banner__afisha > a", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(link)
	})

	c.OnHTML("article[itemscope]", func(e *colly.HTMLElement) {
		data := make(map[string]any)
		// Age restriction
		age := e.ChildText("h1.event-information__name > span.age-limit")
		data["ageRestriction"], _ = strconv.Atoi(strings.Trim(age, "+"))

		// Name
		data["tittle"] = strings.TrimSpace(strings.TrimRight(e.ChildText("h1.event-information__name"), age))

		// Description
		data["description"] = strings.TrimSpace(e.ChildText("div.event-information__description"))

		// ArtworkURL
		data["artworkURL"] = strings.Trim(e.ChildAttr("div.event-information-banner", "style"), "background-image: url()")

		// Genres
		genres := make([]string, 0)
		e.ForEach("div.bottom-tags-wrapper > span.btn-tag-genres", func(i int, e *colly.HTMLElement) {
			genres = append(genres, strings.TrimSpace(e.Text))
		})
		data["genres"] = genres

		// LineUp
		lineup := make([]string, 0)
		e.ForEach(`ul.program-block__list > li.program-block__item[itemprop="performer"]`, func(i int, e *colly.HTMLElement) {
			lineup = append(lineup, squashSpace(e.Text))
		})
		data["lineup"] = lineup
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
		data["startTime"], _ = time.ParseInLocation(
			time.DateTime,
			fmt.Sprintf(
				"%s %s:00",
				e.ChildAttr(`p > time[itemprop="startDate"]`, "datetime"),
				strings.Fields(e.ChildText(`p > time[itemprop="startDate"]`))[3],
			),
			loc,
		)
		// End date. If not presented time.Time zero value
		data["endTime"], _ = time.ParseInLocation(
			time.DateTime,
			fmt.Sprintf(
				"%s %s:00",
				e.ChildAttr(`p > time[itemprop="endDate"]`, "datetime"),
				strings.TrimSpace(e.ChildText(`p > time[itemprop="endDate"]`)),
			),
			loc,
		)

		// Venue
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

		// Tickets
		// TODO: get full url instead of shorten
		data["ticketsURL"] = e.ChildAttr("a.buy-btn", "href")

		// output <- event
	})

	// Create another collector to scrape event details
	// details := p.collector.Clone()
	return c
}
