package web

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"iditusi/internal/models"
	"iditusi/internal/parsers"

	"github.com/gocolly/colly"
)

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0"

type Parser struct {
	name      string
	url       string
	collector *colly.Collector
	output    chan parsers.Result
}

func NewWebPage(name string, URL string, collector *colly.Collector) *Parser {
	return &Parser{
		name:      name,
		url:       URL,
		collector: collector,
		output:    make(chan parsers.Result),
	}
}

func (p *Parser) Parse() <-chan parsers.Result {
	p.collector.Visit(p.url)

	go func() {
		p.collector.Wait()
		// details.Wait()
		close(p.output)
	}()

	return p.output
}

var callbacks = map[string]func(collector *colly.Collector){
	"ptichka": func(collector *colly.Collector) {

	},
}

func NewPtichkaWebPage() *Parser {
	const Name = "ptichka"
	const URL = "https://listing.events/spb"
	const Domain = "listing.events"
	collector := colly.NewCollector(
		colly.UserAgent(UserAgent),
		colly.AllowedDomains(Domain),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./.cache/"+Name),
		colly.Async(true),
		// colly.Debugger(&debug.LogDebugger{}),
		func(collector *colly.Collector) {

			collector.OnRequest(func(r *colly.Request) {
				log.Println("Visiting URL: ", r.URL.String())
			})

			collector.OnHTML("article > div.event-information-banner__afisha > a", func(e *colly.HTMLElement) {
				link := e.Request.AbsoluteURL(e.Attr("href"))
				collector.Visit(link)
			})

			collector.OnHTML("article[itemscope]", func(e *colly.HTMLElement) {
				// var event models.Event
				// var place models.Venue
				var event parsers.Result
				// Age restriction
				age := e.ChildText("h1.event-information__name > span.age-limit")
				event.AgeRestriction, _ = strconv.Atoi(strings.Trim(age, "+"))
				// Name
				event.Name = strings.TrimSpace(strings.TrimRight(e.ChildText("h1.event-information__name"), age))
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
				artists := make([]models.LineUpArtist, 0)
				e.ForEach(`ul.program-block__list > li.program-block__item[itemprop="performer"]`, func(i int, e *colly.HTMLElement) {
					artists = append(artists, models.LineUpArtist{Name: squashSpace(e.Text)})
				})
				event.LineUp = models.LineUp{
					{"main", artists},
				}
				// Dates
				loc, _ := time.LoadLocation("Europe/Moscow")
				// Start date and time
				event.StartDate, _ = time.ParseInLocation(
					time.DateTime,
					fmt.Sprintf(
						"%s %s:00",
						e.ChildAttr(`p > time[itemprop="startDate"]`, "datetime"),
						strings.Fields(e.ChildText(`p > time[itemprop="startDate"]`))[3],
					),
					loc,
				)
				// End date. If not presented time.Time zero value
				event.EndDate, _ = time.ParseInLocation(
					time.DateTime,
					fmt.Sprintf(
						"%s %s:00",
						e.ChildAttr(`p > time[itemprop="endDate"]`, "datetime"),
						strings.TrimSpace(e.ChildText(`p > time[itemprop="endDate"]`)),
					),
					loc,
				)

				// Venue name
				event.Location.Name = strings.TrimSpace(e.ChildText("p.venue-name > a"))
				// Venue address
				address := strings.Split(e.ChildText("p.venue-adress"), "м.") // typo is left intentionally

				event.Location.Address = squashSpace(address[0])

				// Venue metro stations
				stations := make([]string, 0)
				if len(address) > 1 {
					metro := strings.Split(address[1], ",")
					for _, i := range metro {
						stations = append(stations, strings.TrimSpace(i))
					}
				}
				event.Location.NearestMetroStations = stations

				// Tickets
				// TODO: get full url instead of shorten
				event.Tickets.URL = e.ChildAttr("a.buy-btn", "href")

				output <- event
			})
		},
	)

	page := NewWebPage(Name, URL, collector)

	// Create another collector to scrape event details
	// details := p.collector.Clone()

	return page
}

// Locations
// c.OnHTML("article > .item-info", func(e *colly.HTMLElement) {
// 	// fmt.Println(e.Name, e.Attr("class"))
//
// 	venue := location.Venue{
// 		ID:   id.NewID(),
// 		Type: location.Bar,
// 	}
//
// 	// scraping the data of interest
// 	venue.Name = e.ChildText("h3.venue-name a")
// 	div := e.ChildText("div.address")
//
// 	if strings.Contains(div, "\n") {
// 		address := strings.Split(div, "\n")
// 		metroString := regexp.MustCompile(`\s+`).ReplaceAllString(strings.ReplaceAll(address[1], "м.", ""), " ")
// 		metroArray := strings.Split(strings.TrimSpace(metroString), ", ")
//
// 		venue.MetroStations = metroArray
// 		venue.Address = address[0]
// 	}
//
// 	venue.Address = strings.TrimSpace(venue.Address)
//
// 	e.ForEach("div > span.btn-tag-genres", func(_ int, e *colly.HTMLElement) {
// 		venue.MusicGenres = append(venue.MusicGenres, strings.TrimSpace(e.Text))
// 	})
//
// 	// adding the product instance with scraped data to the list of products
//
//
// 	fmt.Println(venue.ID, venue.Name, venue.Address, venue.MusicGenres, venue.MetroStations)
// })
