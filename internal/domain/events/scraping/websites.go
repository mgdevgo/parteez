package scraping

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

func NewRuporEventsWebsite(logger *slog.Logger) Source {
	url, _ := url.Parse("https://rupor.events/spb")

	linkCallback := func(website *Website) colly.CollectorOption {
		return func(c *colly.Collector) {
			c.OnHTML("article > div > a", func(e *colly.HTMLElement) {
				link := e.Request.AbsoluteURL(e.Attr("href"))
				if err := c.Visit(link); err != nil {
					logger.Error("Failed to fetch event details page", "source_id", website.ID(), "url", link, "error", err)
				}
			})
		}
	}

	dataCallback := func(website *Website) colly.CollectorOption {
		return func(c *colly.Collector) {
			c.OnHTML("article[itemscope]", func(e *colly.HTMLElement) {
				errors := make([]error, 0)
				// Age
				age := strings.Trim(e.ChildText("h1.event-information__name > span.age-limit"), "+")
				if age == "" {
					errors = append(errors, fmt.Errorf("age is empty"))
				}
				// Title
				// title := strings.TrimSpace(strings.TrimRight(e.ChildText("h1.event-information__name"), age))
				title := strings.TrimSpace(e.ChildText("h1.event-information__name"))
				if title == "" {
					errors = append(errors, fmt.Errorf("title is empty"))
				}
				// Description
				description := strings.TrimSpace(e.ChildText("div.event-information__description"))
				// ArtworkURL
				artworkURL := strings.Trim(e.ChildAttr("div.event-information-banner", "style"), "background-image: url()")
				// Genres
				genres := make([]string, 0)
				e.ForEach("div.bottom-tags-wrapper > span.btn-tag-genres", func(i int, e *colly.HTMLElement) {
					genres = append(genres, strings.TrimSpace(e.Text))
				})
				// LineUp
				lineup := make([]string, 0)
				e.ForEach(`ul.program-block__list > li.program-block__item[itemprop="performer"]`, func(i int, e *colly.HTMLElement) {
					lineup = append(lineup, squashSpace(e.Text))
				})
				// StartDate
				start := fmt.Sprintf("%s %s:00",
					e.ChildAttr(`p > time[itemprop="startDate"]`, "datetime"),
					strings.Fields(e.ChildText(`p > time[itemprop="startDate"]`))[3],
				)

				endDate := fmt.Sprintf("%s %s:00",
					e.ChildAttr(`p > time[itemprop="endDate"]`, "datetime"),
					strings.TrimSpace(e.ChildText(`p > time[itemprop="endDate"]`)),
				)
				// VenueName
				venueName := strings.TrimSpace(e.ChildText("p.venue-name > a"))
				// VenueAddress
				address := strings.Split(e.ChildText("p.venue-adress"), "м.")
				venueAddress := squashSpace(address[0])
				// MetroStations
				stations := make([]string, 0)
				if len(address) > 1 {
					metro := strings.Split(address[1], ",")
					for _, i := range metro {
						stations = append(stations, strings.TrimSpace(i))
					}
				}

				ticketsURL := e.ChildAttr("a.buy-btn", "href")

				website.result <- Result{
					Event: Event{
						Title:          title,
						Description:    description,
						AgeRestriction: age,
						LineUp:         lineup,
						Genres:         genres,
						StartDate:      start,
						EndDate:        endDate,
						ArtworkURL:     artworkURL,
						TicketsURL:     ticketsURL,
					},
					Venue: Venue{
						Name:          venueName,
						Address:       venueAddress,
						MetroStations: stations,
					},
					Errors: errors,
				}
			})
		}
	}

	return NewWebsite(url, logger, linkCallback, dataCallback)
}

func NewBlankWebsite(logger *slog.Logger) Source {
	url, _ := url.Parse("https://blankclub.ru")

	return NewWebsite(url, logger)

}
